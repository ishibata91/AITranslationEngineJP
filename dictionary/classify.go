package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	wordNetName    = "Japanese WordNet"
	wordNetVersion = "1.1"
)

type classifyResult struct {
	Senses   int            `json:"senses"`
	Matches  int            `json:"matches"`
	Statuses map[string]int `json:"statuses"`
}

func (s *store) classifyGeneralDictionary(ctx context.Context, wordNetPath string) (classifyResult, error) {
	abs, err := filepath.Abs(wordNetPath)
	if err != nil {
		return classifyResult{}, fmt.Errorf("日本語WordNetのpath解決: %w", err)
	}
	if _, statErr := os.Stat(abs); statErr != nil {
		return classifyResult{}, fmt.Errorf("日本語WordNetの確認: %w", statErr)
	}
	if _, attachErr := s.db.ExecContext(ctx, `ATTACH DATABASE ? AS wordnet`, abs); attachErr != nil {
		return classifyResult{}, fmt.Errorf("日本語WordNetの接続: %w", attachErr)
	}
	defer func() {
		_, _ = s.db.ExecContext(context.Background(), `DETACH DATABASE wordnet`)
	}()
	if validationErr := validateWordNet(ctx, s); validationErr != nil {
		return classifyResult{}, validationErr
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return classifyResult{}, fmt.Errorf("一般辞書分類transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM general_dictionary_match
		WHERE dictionary_name = ? AND dictionary_version = ?
		  AND match_status IN ('no_english_headword', 'same_surface_only', 'same_mean_candidate')`,
		wordNetName, wordNetVersion); err != nil {
		return classifyResult{}, fmt.Errorf("前回の自動照合削除: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO general_dictionary_match
		    (sense_id, dictionary_name, dictionary_version, external_sense_id,
		     part_of_speech, definition, japanese_lemma, match_status, reason)
		SELECT DISTINCT s.id, ?, ?, es.synset, ew.pos,
		       COALESCE((
		           SELECT sd.def FROM wordnet.synset_def sd
		           WHERE sd.synset = es.synset AND sd.lang = 'eng'
		           ORDER BY length(sd.def) DESC, sd.def LIMIT 1
		       ), ''),
		       s.dest, 'same_mean_candidate',
		       '原語と訳語が日本語WordNetの同じsynsetにある。Skyrimでの意味は未確認'
		FROM dictionary_sense s
		JOIN dictionary_term t ON t.id = s.term_id
		JOIN wordnet.word ew
		  ON ew.lang = 'eng' AND ew.lemma = lower(replace(t.source, ' ', '_'))
		JOIN wordnet.sense es ON es.wordid = ew.wordid AND es.lang = 'eng'
		JOIN wordnet.sense js ON js.synset = es.synset AND js.lang = 'jpn'
		JOIN wordnet.word jw ON jw.wordid = js.wordid AND jw.lang = 'jpn' AND jw.lemma = s.dest
		WHERE NOT EXISTS (
		    SELECT 1 FROM general_dictionary_match existing
		    WHERE existing.sense_id = s.id
		      AND existing.dictionary_name = ?
		      AND existing.dictionary_version = ?
		      AND existing.external_sense_id = es.synset
		      AND existing.match_status IN ('same_mean_and_translation', 'different_meaning_or_translation')
		)`, wordNetName, wordNetVersion, wordNetName, wordNetVersion); err != nil {
		return classifyResult{}, fmt.Errorf("同じ意味と訳の候補保存: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO general_dictionary_match
		    (sense_id, dictionary_name, dictionary_version, external_sense_id,
		     part_of_speech, definition, match_status, reason)
		SELECT DISTINCT s.id, ?, ?, es.synset, ew.pos,
		       COALESCE((
		           SELECT sd.def FROM wordnet.synset_def sd
		           WHERE sd.synset = es.synset AND sd.lang = 'eng'
		           ORDER BY length(sd.def) DESC, sd.def LIMIT 1
		       ), ''),
		       'same_surface_only',
		       '原語に対応する英語見出しはあるが、同じsynsetに訳語の完全一致がない'
		FROM dictionary_sense s
		JOIN dictionary_term t ON t.id = s.term_id
		JOIN wordnet.word ew
		  ON ew.lang = 'eng' AND ew.lemma = lower(replace(t.source, ' ', '_'))
		JOIN wordnet.sense es ON es.wordid = ew.wordid AND es.lang = 'eng'
		WHERE NOT EXISTS (
		    SELECT 1
		    FROM wordnet.sense js
		    JOIN wordnet.word jw ON jw.wordid = js.wordid AND jw.lang = 'jpn'
		    WHERE js.synset = es.synset AND js.lang = 'jpn' AND jw.lemma = s.dest
		)
		AND NOT EXISTS (
		    SELECT 1 FROM general_dictionary_match existing
		    WHERE existing.sense_id = s.id
		      AND existing.dictionary_name = ?
		      AND existing.dictionary_version = ?
		      AND existing.external_sense_id = es.synset
		      AND existing.match_status IN ('same_mean_and_translation', 'different_meaning_or_translation')
		)`, wordNetName, wordNetVersion, wordNetName, wordNetVersion); err != nil {
		return classifyResult{}, fmt.Errorf("同じ表記だけの候補保存: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO general_dictionary_match
		    (sense_id, dictionary_name, dictionary_version, match_status, reason)
		SELECT s.id, ?, ?, 'no_english_headword',
		       '原語に対応する英語見出しが日本語WordNetにない'
		FROM dictionary_sense s
		JOIN dictionary_term t ON t.id = s.term_id
		WHERE NOT EXISTS (
		    SELECT 1 FROM wordnet.word ew
		    WHERE ew.lang = 'eng' AND ew.lemma = lower(replace(t.source, ' ', '_'))
		)`, wordNetName, wordNetVersion); err != nil {
		return classifyResult{}, fmt.Errorf("英語見出しなしの判定保存: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE dictionary_sense
		SET general_match_status = CASE
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'same_mean_and_translation'
		        ) THEN 'same_mean_and_translation'
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'same_mean_candidate'
		        ) THEN 'same_mean_candidate'
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'same_surface_only'
		        ) THEN 'same_surface_only'
		        WHEN EXISTS (
		            SELECT 1 FROM general_dictionary_match m
		            WHERE m.sense_id = dictionary_sense.id
		              AND m.match_status = 'different_meaning_or_translation'
		        ) THEN 'different_meaning_or_translation'
		        ELSE 'no_english_headword'
		    END,
		    classification_status = CASE
		        WHEN classification_status = 'unclassified' THEN 'general_dictionary_checked'
		        ELSE classification_status
		    END,
		    updated_at = CURRENT_TIMESTAMP`); err != nil {
		return classifyResult{}, fmt.Errorf("意味候補の分類状態更新: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return classifyResult{}, fmt.Errorf("一般辞書分類commit: %w", err)
	}
	return s.classifyStatus(ctx)
}

func validateWordNet(ctx context.Context, s *store) error {
	for _, table := range []string{"word", "sense", "synset_def"} {
		var exists int
		if err := s.db.GetContext(ctx, &exists, `
			SELECT COUNT(*) FROM wordnet.sqlite_schema WHERE type = 'table' AND name = ?`, table); err != nil {
			return fmt.Errorf("日本語WordNetのschema確認: %w", err)
		}
		if exists == 0 {
			return fmt.Errorf("日本語WordNetに必要なtable %q がない", table)
		}
	}
	return nil
}

func (s *store) classifyStatus(ctx context.Context) (classifyResult, error) {
	var out classifyResult
	if err := s.db.GetContext(ctx, &out.Senses, `SELECT COUNT(*) FROM dictionary_sense`); err != nil {
		return classifyResult{}, fmt.Errorf("分類した意味候補数の取得: %w", err)
	}
	if err := s.db.GetContext(ctx, &out.Matches, `SELECT COUNT(*) FROM general_dictionary_match`); err != nil {
		return classifyResult{}, fmt.Errorf("一般辞書照合数の取得: %w", err)
	}
	statuses, err := s.countBy(ctx, "dictionary_sense", "general_match_status")
	if err != nil {
		return classifyResult{}, err
	}
	out.Statuses = statuses
	return out, nil
}

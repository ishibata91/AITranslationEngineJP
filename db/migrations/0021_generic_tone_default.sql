-- 旧既定値を編集せず使っているDBだけ、職業を仮定しない汎用台詞の既定指示へ移行する。
UPDATE tone_default
SET generic_tone_text = '話者を特定できない汎用的な台詞。特定の職業や立場を仮定せず、原文に合う自然な口調で訳す。'
WHERE id = 1
  AND generic_tone_text = '衛兵などの不特定多数が話す汎用的な台詞。職務的で簡潔な口調で訳す。';

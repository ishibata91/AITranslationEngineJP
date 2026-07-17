package tone

// TRDT 由来の台詞感情型。INFO 応答が持つ「その台詞を喋る時の表情演出」の型を表す。
// 値は抽出器（Mutagen の Emotion 型の ToString()）および中心 DB の emotion_type と一致させる文字列にする。
// 対人段階×感情段階の基底口調（9 セル、EmotionSuppressed 等）とは別軸で、台詞ごとにプロンプトへ重ねる。
// Neutral は上乗せしない既定として持つが、DB へは書かない（line_emotion に無い＝加算なし）。
const (
	LineEmotionNeutral  = "Neutral"
	LineEmotionAnger    = "Anger"
	LineEmotionDisgust  = "Disgust"
	LineEmotionFear     = "Fear"
	LineEmotionSad      = "Sad"
	LineEmotionHappy    = "Happy"
	LineEmotionSurprise = "Surprise"
	LineEmotionPuzzled  = "Puzzled"
)

// プロンプトの指示文（directive）の表示用の型。値（派生関数）は directive-presentation.ts に置く。
// プロンプト = Base 指示 + その REC:FIELD に割り当てた指示文（変数を実行時に埋めたもの）。
// 口調・文体・固有名・定型句はすべて「指示文」という 1 つの形に揃える。違いは変数を持つかだけ。

// 指示文に差し込める変数（例: {traits}）。実行時にその実体（話者など）から埋める。
export interface TemplateVariable {
  // 差し込み口のキー（例: {traits}）。
  token: string
  // 差し込まれる内容の説明。
  description: string
}

// 再利用できる指示文。REC:FIELD ごとの「どう訳すか」を表す。複数の REC:FIELD が 1 つを共有する。
export interface Directive {
  // 指示文のキー（例: 説明体・口調・固有名）。
  key: string
  // 指示文の本文（編集できる）。Base 指示へ続けて合成する。
  instruction: string
  // 差し込める変数（無ければ空）。
  variables: TemplateVariable[]
}

// REC:FIELD と、割り当てた指示文キーの対応（固定）。
export interface RecordAssignment {
  // REC:FIELD（例: WEAP:DESC）。
  recField: string
  // 論理名（例: 武器の説明）。コードだけでは分からない意味を添える。
  logicalName: string
  // 割り当てた指示文のキー。
  directive: string
}

// 指示文と、その指示文を割り当てた REC:FIELD（対象）を束ねた表示用の派生値。
export interface DirectiveSection {
  key: string
  instruction: string
  variables: TemplateVariable[]
  // この指示文を割り当てた REC:FIELD。読み取り専用で出す。
  targets: RecordAssignment[]
}

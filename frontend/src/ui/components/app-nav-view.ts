// 画面間ナビゲーションの表示用の型。遷移ロジックは持たず、現在画面と遷移先だけを表す。

// アプリの画面。plugins=翻訳対象プラグイン、translation=翻訳実行、template=プロンプトテンプレート編集。
export type AppRoute = "plugins" | "translation" | "template" | "termDictionary"

// ナビゲーション項目。route=遷移先の画面、label=タブに出す表示名。
export interface NavItem {
  route: AppRoute
  label: string
}

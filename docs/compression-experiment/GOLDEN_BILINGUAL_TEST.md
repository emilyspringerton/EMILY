# 黃金上下文 — 2026-06-12
*雙語壓縮測試版 / Bilingual compression test. Source: EMILY/BACKLOG.md.*

## 系統 (SYSTEM)
倉庫: EMILY · FATBABY · IDUNA · TYLER · SHANKPIT · MJOLNIR · APPLES · emily.cli
接口: IDUNA /api/v1/heimdal/sprints | 北極星: FATBABY/docs/northstar/northstar.md
積壓: EMILY/BACKLOG.md | 完成: EMILY/DONE.md

## 積壓狀態 (BACKLOG): 開放24 完成10 待分10

## 已完成區段 (no open items): S6:RSI · S12:HEIMDAL

## 開放事項 (OPEN)

### S2 影片管道 MPT (3待)
• Pexels密鑰 — config.toml YOUR_PEXELS_API_KEY_HERE 等待人工操作
• S01E01冷開場片段 — 依賴: Pexels密鑰
• MPT→TYLER觸發 — 新episode腳本→自動調用MPT

### S3 系統可見性 (2待)
• Apple儀器審計 — cron觸發/觀察遺漏/haiku失敗/HEIMDAL未記錄
• 單一日誌流 — obs-watcher/rsi-loop/emily-agent/IDUNA→ndjson→git同步

### S4 遊戲引擎 (3待)
• SHANKPIT→MPT橋接 — 等待MPT端到端運行
• TYLER手機角色擴展 — 通知/應用/聯繫人/地圖規格
• TYLER過場對話系統 — FFXI風格對話場景

### S5 未來 (2待)
• MySQL嵌入式服務器 — go-mysql-server，依賴穩定SQLite
• Tyler IDUNA代理CLI註冊 — 等待iduna CLI構建

### S10 網頁審計 (1待)
• 前門驗證器 — MJOLNIR WebView目標上線後運行web_audit_url

### S13 黃金文件 (2待)
• 文件審計 — 所有倉庫northstar/golden文件審查
• Emily Prime API一致性 — emily-agent暴露emily.cli命令集

### S14 EmilyOS (2待)
• EmilyOS北極星 — 裸機外核northstar文件
• 包倉庫+構建系統 — Debian/Arch基礎，可重現構建

### S15 PITVIPER (1待)
• PITVIPER北極星 — SDL2+FreeType2終端northstar

### S16 AI層 FABLE (3待)
• FABLE顧問(基本) — haiku→FABLE顧問工具
• FABLE→HEIMDAL整合 — 推薦→sprint自動排隊
• Emily Prime API — 外部編排穩定API接口

### S17 新聞站+GTM (5待)
• 股票圖表 — SVG內聯股價圖表
• 500錯誤調查 — 新聞站500s根因
• Emily原創評論接口 — Emily Prime POST文章入CMS
• GTM漏斗 — Ask Emily免費→訂閱→社區→Merkle查詢
• 自我改進訓練管道 — 用戶數據飛輪，RLHF循環

### 待分類 (10項 — 運行 emily backlog promote)
• UX: ticker搜索點擊/Enter自動跳轉 — 去除冗餘Go按鈕
• 功能: Emily觀察時自動創建GitHub issue
• 所有env變量記錄在README頂部
• obs-watcher向Claude Code提示注入完整報告+git同步要求
• EDGAR submissions截斷JSON — BAC/C/GS/JPM/MS
• entity-graph 8-K文件檢測 — form/source_type不匹配
• entity-graph 8-K子類型解析 — Item 5.07未找到
• entity-graph讀取0歸檔 — 846源文件但游標問題
• ...還有2項

## RSI核心指令
管道: 原始想法→emily eo→待分類→emily backlog promote(haiku)→編號區段→HEIMDAL sprint→實施
每次完成: Apple(IDUNA POST /api/v1/apples)+標記[x]+git commit+push

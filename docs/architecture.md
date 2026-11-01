# cardVault — Architecture

## Komponen
- Wails app (Go) menyediakan window webview + service bindings.
- Frontend SvelteKit (SPA, ssr=false) berkomunikasi via Wails runtime Call.
- SQLite single file (modernc, pure Go) untuk mobile compatibility.

## Data Flow
UI (Svelte) → Wails binding → Service (Go) → Repository → SQLite → response.

## Services (Wails-bound)
- CardService: ListCards, GetCard, AddCard, UpdateCard, DeleteCard, TouchCard.
- ProviderService: ListProviders, GetProvider, AddProvider, DeleteProvider.
- NFCService: Capture, ListCaptures, Simulate.

## Mobile
- Android/iOS dibangun via `wails3 android/ios overlay:gen` yang menghasilkan project
  native; Go core (services + SQLite) direuse.
- NFC hardware diakses dari layer native, memanggil service Go lewat bridge.

## Storage
- Default: ~/.cardvault/cardvault.db (desktop).
- Mobile: app sandbox path (set via CARDVAULT_DB env atau config).
- Backup: export/import file SQLite.
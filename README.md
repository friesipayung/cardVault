# cardVault

Card loyalty management app built with Wails v3 (Go) + SvelteKit (Svelte 5) + SQLite.

## Fitur
- Vault kartu loyalty (provider preset / custom).
- Tambah kartu via nomor/nama/identifier.
- Scan NFC card → simpan info; simulasi NFC saat dipakai.
- Database SQLite (single file).

## Dev
```bash
export PATH="$HOME/go/bin:/home/friendly/.hermes/node/bin:$PATH"
wails3 dev
```

## Build desktop
```bash
wails3 build
```

## Build mobile
```bash
wails3 android overlay:gen
wails3 ios xcode:gen
```

Lihat docs/plan.md dan docs/architecture.md.
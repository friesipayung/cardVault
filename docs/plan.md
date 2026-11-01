# cardVault — Plan

## Goal
Card loyalty management app. Store loyalty cards (preset providers + custom), scan card codes, read NFC cards and simulate (HCE) them.

## Stack
- Wails v3 (Go backend + web frontend)
- SvelteKit + Svelte 5 (frontend)
- Go
- SQLite (binary, file-based)
- Android (arm64-v8a + x86_64), iOS later

## Features
1. Card vault: add/edit/delete loyalty cards (provider, number, name, barcode/QR).
2. Preset providers: quick add by number/name/identifier.
3. Barcode/QR scan (camera) for card codes.
4. NFC: read NFC card → save; simulate card via HCE when paying.

## Android build
- Package: com.cardvault.app
- minSdk 23 (Android 6), targetSdk 35
- ABIs: arm64-v8a, x86_64
- 16 KB page aligned (Android 15)
- CI: GitHub Actions builds APK, tests install on emulator, uploads to release.

## NFC (planned)
- Permission: android.permission.NFC
- Feature: android.hardware.nfc (optional)
- HCE service for card emulation (requires BIND_NFC_SERVICE).

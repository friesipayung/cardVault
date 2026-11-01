<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import QRCode from 'qrcode';
  import JsBarcode from 'jsbarcode';
  import * as CardService from '../../bindings/github.com/friesipayung/cardVault/cardservice.js';
  import * as ProviderService from '../../bindings/github.com/friesipayung/cardVault/providerservice.js';
  import * as NFCService from '../../bindings/github.com/friesipayung/cardVault/nfcservice.js';
  import type { Card } from '../../bindings/github.com/friesipayung/cardVault/internal/card/models.js';
  import type { Provider } from '../../bindings/github.com/friesipayung/cardVault/internal/provider/models.js';

  let cards: Card[] = $state([]);
  let providers: Provider[] = $state([]);
  let nfcUid = $state('');
  let nfcCardId = $state<number | null>(null);
  let nfcResult = $state('');
  let search = $state('');
  let form = $state({ provider_name: '', card_number: '', card_holder_name: '', barcode_type: 'QR', barcode_value: '' });
  let providerForm = $state({ name: '', barcode_format: 'Barcode', country: '' });
  let loading = $state(false);
  let scanning = $state(false);
  let scanner: any = null;

  async function load() {
    loading = true;
    try {
      cards = (await CardService.ListCards(search)) ?? [];
      providers = (await ProviderService.ListProviders('')) ?? [];
    } finally {
      loading = false;
    }
  }

  async function add() {
    if (!form.provider_name || !form.card_number) return;
    await CardService.AddCard({
      id: 0,
      provider_id: null,
      provider_name: form.provider_name,
      card_number: form.card_number,
      card_holder_name: form.card_holder_name,
      barcode_type: form.barcode_type,
      barcode_value: form.barcode_value || form.card_number,
      nfc_uid: '',
      notes: '',
      created_at: '',
      updated_at: '',
      last_used_at: null,
    });
    form = { provider_name: '', card_number: '', card_holder_name: '', barcode_type: 'QR', barcode_value: '' };
    await load();
  }

  async function remove(id: number) {
    await CardService.DeleteCard(id);
    await load();
  }

  async function addProvider() {
    if (!providerForm.name) return;
    await ProviderService.AddProvider({
      id: 0,
      name: providerForm.name,
      logo_url: '',
      barcode_format: providerForm.barcode_format,
      country: providerForm.country,
      notes: '',
    });
    providerForm = { name: '', barcode_format: 'Barcode', country: '' };
    await load();
  }

  async function removeProvider(id: number) {
    await ProviderService.DeleteProvider(id);
    await load();
  }

  async function captureNfc() {
    if (!nfcUid) return;
    const res = await NFCService.Capture(nfcUid, '', '', '', nfcCardId);
    nfcResult = `Disimpan: ${res.uid}`;
  }

  async function simulateNfc() {
    if (!nfcCardId) return;
    const uid = await NFCService.Simulate(nfcCardId);
    nfcResult = `Simulasi UID: ${uid}`;
  }

  function renderCode(node: HTMLCanvasElement, card: Card) {
    const draw = () => {
      const value = card.barcode_value || card.card_number;
      if (card.barcode_type === 'QR') {
        QRCode.toCanvas(node, value, { width: 128, margin: 1 });
      } else {
        JsBarcode(node, value, { format: 'CODE128', width: 2, height: 60, displayValue: false });
      }
    };
    draw();
    return {
      update(c: Card) {
        card = c;
        draw();
      },
    };
  }

  async function startScan() {
    scanning = true;
    const { Html5Qrcode } = await import('html5-qrcode');
    scanner = new Html5Qrcode('reader');
    await scanner.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: 220 },
      (decoded: string) => {
        form.card_number = decoded;
        form.barcode_value = decoded;
        stopScan();
      },
      () => {}
    );
  }

  async function stopScan() {
    if (scanner) {
      try { await scanner.stop(); } catch {}
      scanner = null;
    }
    scanning = false;
  }

  onMount(load);
  onDestroy(() => { stopScan(); });
</script>

<main>
  <h1>cardVault</h1>
  <p>Loyalty card vault (Wails v3 + SvelteKit + Svelte 5)</p>

  <section class="card">
    <h2>Tambah Kartu</h2>
    <label>Provider
      <select bind:value={form.provider_name}>
        <option value="">-- pilih --</option>
        {#each providers as p}
          <option value={p.name}>{p.name}</option>
        {/each}
      </select>
    </label>
    <label>Nomor Kartu <input bind:value={form.card_number} /></label>
    <label>Nama Pemegang <input bind:value={form.card_holder_name} /></label>
    <label>Tipe Barcode
      <select bind:value={form.barcode_type}>
        <option>QR</option>
        <option>Barcode</option>
      </select>
    </label>
    <label>Nilai Barcode (opsional) <input bind:value={form.barcode_value} /></label>
    <div class="row">
      <button onclick={add}>Simpan</button>
      <button onclick={startScan} disabled={scanning}>Scan Kamera</button>
      {#if scanning}<button onclick={stopScan}>Stop</button>{/if}
    </div>
    {#if scanning}
      <div id="reader" class="reader"></div>
    {/if}
  </section>

  <section class="card">
    <h2>Kartu Tersimpan</h2>
    <input placeholder="cari..." bind:value={search} oninput={load} />
    {#if loading}<p>Memuat...</p>{/if}
    {#each cards as c (c.id)}
      <div class="card row-card">
        <div class="code">
          <canvas use:renderCode={c}></canvas>
        </div>
        <div class="info">
          <strong>{c.provider_name}</strong> — {c.card_number}
          <div>{c.card_holder_name}</div>
          <div class="row">
            <button onclick={() => remove(c.id)}>Hapus</button>
          </div>
        </div>
      </div>
    {/each}
    {#if cards.length === 0 && !loading}<p>Belum ada kartu.</p>{/if}
  </section>

  <section class="card">
    <h2>Provider</h2>
    <div class="row">
      <input placeholder="nama provider" bind:value={providerForm.name} />
      <select bind:value={providerForm.barcode_format}>
        <option>QR</option>
        <option>Barcode</option>
      </select>
      <input placeholder="negara" bind:value={providerForm.country} />
      <button onclick={addProvider}>Tambah</button>
    </div>
    <ul>
      {#each providers as p}
        <li>{p.name} ({p.barcode_format}) <button onclick={() => removeProvider(p.id)}>x</button></li>
      {/each}
    </ul>
  </section>

  <section class="card">
    <h2>NFC</h2>
    <label>UID <input bind:value={nfcUid} placeholder="scan UID" /></label>
    <label>Tautkan ke Kartu
      <select bind:value={nfcCardId}>
        <option value={null}>-- tidak ditautkan --</option>
        {#each cards as c}
          <option value={c.id}>{c.provider_name} {c.card_number}</option>
        {/each}
      </select>
    </label>
    <div class="row">
      <button onclick={captureNfc}>Simpan Capture</button>
      <button onclick={simulateNfc}>Simulasi NFC</button>
    </div>
    {#if nfcResult}<p>{nfcResult}</p>{/if}
  </section>
</main>

<style>
  main { max-width: 820px; margin: 0 auto; padding: 24px; }
  label { display: block; margin: 8px 0; }
  input, select { width: 100%; padding: 8px; margin-top: 4px; }
  .row { margin-top: 8px; display: flex; gap: 8px; align-items: center; }
  .row-card { display: flex; gap: 16px; align-items: center; }
  .code { flex: 0 0 140px; text-align: center; }
  .info { flex: 1; }
  ul { padding-left: 20px; }
  .reader { margin-top: 12px; width: 100%; max-width: 320px; }
</style>
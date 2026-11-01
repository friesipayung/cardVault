<script lang="ts">
  import { onMount } from 'svelte';
  import * as CardService from '../../bindings/github.com/friesipayung/cardVault/cardservice.js';
  import * as ProviderService from '../../bindings/github.com/friesipayung/cardVault/providerservice.js';
  import * as NFCService from '../../bindings/github.com/friesipayung/cardVault/nfcservice.js';
  import type { Card } from '../../bindings/github.com/friesipayung/cardVault/internal/card/models.js';

  let cards: Card[] = $state([]);
  let providers: { id: number; name: string }[] = $state([]);
  let nfcUid = $state('');
  let nfcCardId = $state<number | null>(null);
  let nfcResult = $state('');
  let search = $state('');
  let form = $state({ provider_name: '', card_number: '', card_holder_name: '', barcode_type: 'QR', barcode_value: '' });
  let loading = $state(false);

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

  onMount(load);
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
        <option>EAN13</option>
      </select>
    </label>
    <button onclick={add}>Simpan</button>
  </section>

  <section class="card">
    <h2>Kartu Tersimpan</h2>
    <input placeholder="cari..." bind:value={search} oninput={load} />
    {#if loading}<p>Memuat...</p>{/if}
    {#each cards as c (c.id)}
      <div class="card">
        <strong>{c.provider_name}</strong> — {c.card_number}
        <div>{c.card_holder_name}</div>
        <div class="row">
          <button onclick={() => remove(c.id)}>Hapus</button>
        </div>
      </div>
    {/each}
    {#if cards.length === 0 && !loading}<p>Belum ada kartu.</p>{/if}
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
  main { max-width: 720px; margin: 0 auto; padding: 24px; }
  label { display: block; margin: 8px 0; }
  input, select { width: 100%; padding: 8px; margin-top: 4px; }
  .row { margin-top: 8px; display: flex; gap: 8px; }
</style>
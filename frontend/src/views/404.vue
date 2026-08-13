<template>
  <section class="page-404 bv-empty">
    <div class="bv-empty__art" aria-hidden="true">
      <svg width="88" height="72" viewBox="0 0 88 72" fill="none">
        <rect x="10" y="16" width="68" height="46" rx="10" fill="#F4F7F5" stroke="#E5E7EB" />
        <text x="44" y="48" text-anchor="middle" font-size="22" font-weight="700" fill="#111">404</text>
      </svg>
    </div>
    <h1 class="bv-empty__title">Page not found</h1>
    <p class="bv-empty__lead">This page does not exist or you do not have access to it.</p>
    <div class="bv-empty__actions">
      <router-link class="btn-new" :to="{ name: 'dashboard' }">Back to Home</router-link>
    </div>
  </section>
</template>

<script>
import Vue from 'vue';

export default Vue.extend({});
</script>

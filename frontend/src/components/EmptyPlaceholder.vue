<template>
  <section class="bv-empty">
    <div class="bv-empty__art" aria-hidden="true">
      <svg width="88" height="72" viewBox="0 0 88 72" fill="none">
        <rect x="10" y="16" width="68" height="46" rx="10" fill="#F4F7F5" stroke="#E5E7EB" />
        <path d="M24 36h40M24 46h24" stroke="#C5CDD6" stroke-width="2.4" stroke-linecap="round" />
        <circle cx="66" cy="22" r="11" fill="#EAF7EC" stroke="#0B996E" stroke-width="1.6" />
        <path d="M66 17.5v9M61.5 22h9" stroke="#0B996E" stroke-width="1.7" stroke-linecap="round" />
      </svg>
    </div>
    <h2 class="bv-empty__title">{{ titleText }}</h2>
    <p class="bv-empty__lead">{{ descriptionText }}</p>
    <div v-if="$slots.default" class="bv-empty__actions">
      <slot />
    </div>
  </section>
</template>

<script>
export default {
  name: 'EmptyPlaceholder',

  props: {
    icon: { type: String, default: '' },
    label: { type: String, default: '' },
    description: { type: String, default: '' },
  },

  computed: {
    titleText() {
      return this.label || this.$t('globals.messages.emptyState');
    },
    descriptionText() {
      return this.description || 'When you add items, they will appear here.';
    },
  },
};
</script>

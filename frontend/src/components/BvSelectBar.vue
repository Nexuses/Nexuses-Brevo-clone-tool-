<template>
  <div v-if="selected > 0" class="bv-select-bar">
    <div class="bv-select-bar__left">
      <label class="bv-select-bar__check">
        <input type="checkbox" checked aria-label="Clear selection" @change="$emit('clear')" />
      </label>
      <slot />
      <b-dropdown
        v-if="$slots.more"
        class="bv-select-bar__more"
        :mobile-modal="false"
        position="is-bottom-left"
      >
        <template #trigger>
          <button type="button" class="bv-select-bar__action">
            More actions
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5"
                stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
        </template>
        <slot name="more" />
      </b-dropdown>
    </div>
    <div class="bv-select-bar__right">
      <strong>{{ countLabel }}</strong>
      <button
        v-if="showSelectAll"
        type="button"
        class="bv-select-bar__all"
        @click="$emit('select-all')"
      >
        Select all {{ noun }}
      </button>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';

export default Vue.extend({
  name: 'BvSelectBar',

  props: {
    selected: {
      type: Number,
      default: 0,
    },
    total: {
      type: Number,
      default: 0,
    },
    noun: {
      type: String,
      default: 'items',
    },
    allSelected: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    showSelectAll() {
      return !this.allSelected && this.total > this.selected;
    },

    countLabel() {
      const n = this.$utils.formatNumber(this.selected || 0);
      return `${n} ${this.noun} selected`;
    },
  },
});
</script>

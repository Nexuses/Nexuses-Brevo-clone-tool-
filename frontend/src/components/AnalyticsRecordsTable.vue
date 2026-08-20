<template>
  <section>
    <b-table :data="records" :loading="loading" paginated backend-pagination pagination-position="both"
      :total="total" :per-page="perPage" :current-page="page" @page-change="onPageChange" hoverable striped>
      <b-table-column v-if="campaigns.length > 1" field="campaignName" :label="$t('globals.terms.campaign')"
        v-slot="props">
        {{ props.row.campaignName }}
      </b-table-column>
      <b-table-column field="name" :label="$t('globals.fields.name')" v-slot="props">
        <template v-if="props.row.subscriberId">
          <router-link :to="{ name: 'subscriber', params: { id: props.row.subscriberId } }">
            {{ props.row.name || '—' }}
          </router-link>
        </template>
        <span v-else class="has-text-grey-light">—</span>
      </b-table-column>
      <b-table-column field="email" :label="$t('subscribers.email')" v-slot="props">
        <template v-if="props.row.subscriberId">
          {{ props.row.email }}
        </template>
        <span v-else class="has-text-grey-light">{{ $t('analytics.anonymous') }}</span>
      </b-table-column>
      <b-table-column v-if="typ === 'clicks'" field="url" label="Link clicked" v-slot="props">
        <a v-if="props.row.url" :href="props.row.url" target="_blank" rel="noopener noreferrer">
          {{ truncateUrl(props.row.url) }}
        </a>
        <span v-else class="has-text-grey-light">—</span>
      </b-table-column>
      <b-table-column field="createdAt" :label="$t('globals.fields.createdAt')" v-slot="props">
        {{ $utils.duration(new Date(), props.row.createdAt, true) }}
        <span class="has-text-grey-light is-size-7">({{ formatTime(props.row.createdAt) }})</span>
      </b-table-column>

      <template #empty>
        <p class="has-text-grey-light p-4">{{ $t('analytics.noRecords') }}</p>
      </template>
    </b-table>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';

export default Vue.extend({
  props: {
    typ: { type: String, required: true },
    campaigns: { type: Array, default: () => [] },
    from: { type: Date, required: true },
    to: { type: Date, required: true },
    anonymous: { type: Boolean, default: false },
  },

  data() {
    return {
      records: [],
      total: 0,
      page: 1,
      perPage: 30,
      loading: false,
    };
  },

  watch: {
    campaigns: { handler: 'load', deep: true },
    from: 'load',
    to: 'load',
  },

  mounted() {
    this.load();
  },

  methods: {
    formatTime(t) {
      return dayjs(t).format('YYYY-MM-DD HH:mm');
    },

    truncateUrl(url) {
      if (!url || url.length <= 72) {
        return url;
      }
      try {
        const u = new URL(url);
        return `${`${u.hostname}${u.pathname}`.substr(0, 68)}…`;
      } catch {
        return `${url.substr(0, 68)}…`;
      }
    },

    onPageChange(p) {
      this.page = p;
      this.load();
    },

    load() {
      if (!this.campaigns.length) {
        return;
      }
      this.loading = true;
      this.$api.getCampaignAnalyticsRecords(this.typ, {
        id: this.campaigns.map((c) => c.id),
        from: this.from,
        to: this.to,
        page: this.page,
        per_page: this.perPage,
      }).then((data) => {
        this.records = data.results || [];
        this.total = data.total || 0;
        this.loading = false;
      });
    },
  },
});
</script>

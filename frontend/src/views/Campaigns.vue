<template>
  <section class="campaigns campaigns-brevo">
    <marketing-subnav />

    <header class="campaigns-brevo__header">
      <div class="campaigns-brevo__title-row">
        <h1 class="campaigns-brevo__title">Campaigns</h1>
        <div class="campaigns-brevo__actions">
          <button type="button" class="campaigns-brevo__btn campaigns-brevo__btn--outline">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path d="M2.5 4.5h4l1.2 1.5H13.5v6.5a1 1 0 01-1 1h-10a1 1 0 01-1-1V5.5a1 1 0 011-1z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round" />
            </svg>
            Create folder
          </button>
          <button type="button" class="campaigns-brevo__btn campaigns-brevo__btn--ai">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path
                d="M8 2.5l.7 2.1L10.8 5.3 8.7 6.4 8 8.5l-.7-2.1L5.2 5.3l2.1-.7L8 2.5z"
                fill="currentColor"
              />
              <path
                d="M12.5 8.5l.45 1.35 1.35.45-1.35.45-.45 1.35-.45-1.35L10.7 10.3l1.35-.45.45-1.35z"
                fill="currentColor"
              />
              <path
                d="M4.2 9.8l.35 1.05L5.6 11.2l-1.05.35-.35 1.05-.35-1.05L2.8 11.2l1.05-.35.35-1.05z"
                fill="currentColor"
              />
            </svg>
            Generate campaign with AI
          </button>
          <router-link
            v-if="$can('campaigns:manage')"
            :to="{ name: 'campaign', params: { id: 'new' } }"
            class="campaigns-brevo__btn campaigns-brevo__btn--primary"
            data-cy="btn-new"
          >
            Create campaign
          </router-link>
        </div>
      </div>

      <div class="campaigns-brevo__channel-tabs" role="tablist" aria-label="Campaign type">
        <button type="button" class="campaigns-brevo__channel is-active" role="tab" aria-selected="true">
          Email
        </button>
      </div>
    </header>

    <div class="campaigns-brevo__controls">
      <div class="campaigns-brevo__toolbar">
        <form @submit.prevent="getCampaigns" class="campaigns-brevo__search">
          <b-field>
            <b-input
              v-model="queryParams.query"
              name="query"
              expanded
              placeholder="Search for a campaign"
              icon="magnify"
              ref="query"
              @input="onSearchInput"
            />
          </b-field>
        </form>
        <button type="button" class="campaigns-brevo__chip">
          All statuses
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
        <button type="button" class="campaigns-brevo__chip">
          Select tags
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
      </div>
    </div>

    <div class="campaigns-brevo__meta" v-if="campaigns.total > 0 || bulk.checked.length > 0">
      <div class="actions campaigns-brevo__bulk" v-if="bulk.checked.length > 0">
        <a class="a" href="#" @click.prevent="deleteCampaigns" data-cy="btn-delete-campaigns">
          <b-icon icon="trash-can-outline" size="is-small" /> Delete
        </a>
        <span class="a">
          {{ $tc('globals.messages.numSelected', numSelectedCampaigns, { num: numSelectedCampaigns }) }}
        </span>
      </div>
      <div class="campaigns-brevo__pager" v-if="campaigns.total > 0">
        <span class="campaigns-brevo__pager-range">{{ pageRangeLabel }}</span>
        <span class="campaigns-brevo__pager-pages">
          <span class="campaigns-brevo__pager-current">{{ queryParams.page }}</span>
          of {{ totalPages }} pages
        </span>
        <button
          type="button"
          class="campaigns-brevo__pager-btn"
          :disabled="queryParams.page <= 1"
          aria-label="Previous page"
          @click="onPageChange(queryParams.page - 1)"
        >
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
            <path d="M8.5 3.5L5 7l3.5 3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
        <button
          type="button"
          class="campaigns-brevo__pager-btn"
          :disabled="queryParams.page >= totalPages"
          aria-label="Next page"
          @click="onPageChange(queryParams.page + 1)"
        >
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
            <path d="M5.5 3.5L9 7l-3.5 3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
      </div>
    </div>

    <b-table
      class="bv-campaigns-table"
      :data="campaigns.results"
      :loading="loading.campaigns"
      :row-class="highlightedRow"
      :show-header="false"
      @check-all="onTableCheck"
      @check="onTableCheck"
      :checked-rows.sync="bulk.checked"
      paginated
      backend-pagination
      pagination-position="bottom"
      @page-change="onPageChange"
      :current-page="queryParams.page"
      :per-page="campaigns.perPage"
      :total="campaigns.total"
      checkable
      backend-sorting
      @sort="onSort"
    >
      <b-table-column v-slot="props" cell-class="bv-campaign-cell">
        <div class="bv-campaign-card" :set="stats = getCampaignStats(props.row)">
          <div class="bv-campaign-card__main">
            <router-link
              class="bv-campaign-card__title"
              :to="{ name: 'campaign', params: { id: props.row.id } }"
            >
              {{ props.row.name }}
            </router-link>
            <div class="bv-campaign-card__status" :class="`is-${props.row.status}`">
              <span class="dot" />
              <span class="bv-campaign-card__status-label">{{ statusLabel(props.row) }}</span>
              <span class="bv-campaign-card__status-when">{{ statusWhen(props.row, stats) }}</span>
              <span class="spinner is-tiny" v-if="isRunning(props.row.id)">
                <b-loading :is-full-page="false" active />
              </span>
            </div>
            <div class="bv-campaign-card__id">#{{ props.row.id }}</div>
          </div>

          <div class="bv-campaign-card__metrics">
            <div class="bv-metric">
              <span class="bv-metric__lbl">Recipients</span>
              <strong class="bv-metric__val">{{ $utils.formatNumber(recipientsCount(stats)) }}</strong>
              <span class="bv-metric__pct">{{ pctLabel(recipientsCount(stats), recipientsCount(stats)) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Opens</span>
              <strong class="bv-metric__val">{{ $utils.formatNumber(metricVal(stats, ['views', 'opens', 'uniqueOpens', 'unique_opens'])) }}</strong>
              <span class="bv-metric__pct">{{ pctLabel(metricVal(stats, ['views', 'opens', 'uniqueOpens', 'unique_opens']), sentCount(stats)) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Clicks</span>
              <strong class="bv-metric__val">{{ $utils.formatNumber(metricVal(stats, ['clicks'])) }}</strong>
              <span class="bv-metric__pct">{{ pctLabel(metricVal(stats, ['clicks']), sentCount(stats)) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Unsubscribed</span>
              <strong class="bv-metric__val">0</strong>
              <span class="bv-metric__pct">0%</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Conversions</span>
              <strong class="bv-metric__val">—</strong>
              <span class="bv-metric__pct">0%</span>
            </div>
          </div>

          <div class="bv-campaign-card__actions">
            <router-link
              class="bv-campaign-card__report"
              :to="{ name: 'campaign', params: { id: props.row.id }, query: { view: 'report' } }"
              data-cy="btn-campaign-view-report"
              aria-label="Report"
            >
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
                <path d="M3.5 12.5V9.5M7.5 12.5V6.5M11.5 12.5V8M14.5 12.5V4.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
                <path d="M3 14.5h12" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
              </svg>
            </router-link>

            <b-dropdown class="campaign-actions-menu" position="is-bottom-left">
              <template #trigger>
                <button
                  type="button"
                  class="campaign-actions-trigger"
                  data-cy="btn-campaign-actions"
                  aria-label="Actions"
                >
                  <span class="campaign-kebab" aria-hidden="true">
                    <span /><span /><span />
                  </span>
                </button>
              </template>

              <div class="campaign-actions-panel">
                <template v-if="$can('campaigns:send')">
                  <a
                    v-if="canStart(props.row)"
                    href="#"
                    class="campaign-action"
                    @click.prevent="$utils.confirm(null, () => changeCampaignStatus(props.row, 'running'))"
                    data-cy="btn-start"
                    :aria-label="$t('campaigns.start')"
                  >
                    <b-tooltip :label="$t('campaigns.start')" type="is-dark" position="is-left">
                      <b-icon icon="rocket-launch-outline" />
                    </b-tooltip>
                  </a>
                  <a
                    v-else-if="canResume(props.row)"
                    href="#"
                    class="campaign-action"
                    @click.prevent="$utils.confirm(null, () => changeCampaignStatus(props.row, 'running'))"
                    data-cy="btn-resume"
                    :aria-label="$t('campaigns.send')"
                  >
                    <b-tooltip :label="$t('campaigns.send')" type="is-dark" position="is-left">
                      <b-icon icon="rocket-launch-outline" />
                    </b-tooltip>
                  </a>
                  <a
                    v-else-if="canSchedule(props.row)"
                    href="#"
                    class="campaign-action"
                    @click.prevent="$utils.confirm($t('campaigns.confirmSchedule'), () => changeCampaignStatus(props.row, 'scheduled'))"
                    data-cy="btn-schedule"
                    :aria-label="$t('campaigns.schedule')"
                  >
                    <b-tooltip :label="$t('campaigns.schedule')" type="is-dark" position="is-left">
                      <b-icon icon="rocket-launch-outline" />
                    </b-tooltip>
                  </a>
                  <span v-else class="campaign-action is-disabled">
                    <b-icon icon="rocket-launch-outline" />
                  </span>

                  <a
                    v-if="canPause(props.row)"
                    href="#"
                    class="campaign-action"
                    @click.prevent="$utils.confirm(null, () => changeCampaignStatus(props.row, 'paused'))"
                    data-cy="btn-pause"
                    :aria-label="$t('campaigns.pause')"
                  >
                    <b-tooltip :label="$t('campaigns.pause')" type="is-dark" position="is-left">
                      <b-icon icon="pause-circle-outline" />
                    </b-tooltip>
                  </a>

                  <a
                    v-if="canCancel(props.row)"
                    href="#"
                    class="campaign-action"
                    @click.prevent="$utils.confirm(null, () => changeCampaignStatus(props.row, 'cancelled'))"
                    data-cy="btn-cancel"
                    :aria-label="$t('globals.buttons.cancel')"
                  >
                    <b-tooltip :label="$t('globals.buttons.cancel')" type="is-dark" position="is-left">
                      <b-icon icon="cancel" />
                    </b-tooltip>
                  </a>
                  <span v-else class="campaign-action is-disabled">
                    <b-icon icon="cancel" />
                  </span>
                </template>

                <a
                  href="#"
                  class="campaign-action"
                  @click.prevent="previewCampaign(props.row)"
                  data-cy="btn-preview"
                  :aria-label="$t('campaigns.preview')"
                >
                  <b-tooltip :label="$t('campaigns.preview')" type="is-dark" position="is-left">
                    <b-icon icon="file-find-outline" />
                  </b-tooltip>
                </a>

                <a
                  v-if="$can('campaigns:manage')"
                  href="#"
                  class="campaign-action"
                  @click.prevent="$utils.prompt($t('globals.buttons.clone'),
                    {
                      placeholder: $t('globals.fields.name'),
                      value: $t('campaigns.copyOf', { name: props.row.name }),
                    },
                    (name) => cloneCampaign(name, props.row))"
                  data-cy="btn-clone"
                  :aria-label="$t('globals.buttons.clone')"
                >
                  <b-tooltip :label="$t('globals.buttons.clone')" type="is-dark" position="is-left">
                    <b-icon icon="file-multiple-outline" />
                  </b-tooltip>
                </a>

                <a
                  v-if="$can('campaigns:manage')"
                  href="#"
                  class="campaign-action"
                  @click.prevent="$utils.confirm($t('campaigns.confirmDelete', { name: props.row.name }), () => deleteCampaign(props.row))"
                  data-cy="btn-delete"
                  :aria-label="$t('globals.buttons.delete')"
                >
                  <b-tooltip :label="$t('globals.buttons.delete')" type="is-dark" position="is-left">
                    <b-icon icon="trash-can-outline" />
                  </b-tooltip>
                </a>
              </div>
            </b-dropdown>
          </div>
        </div>
      </b-table-column>

      <template #empty v-if="!loading.campaigns">
        <empty-placeholder />
      </template>
    </b-table>

    <campaign-preview
      v-if="previewItem"
      type="campaign"
      :id="previewItem.id"
      :title="previewItem.name"
      @close="closePreview"
    />
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';
import CampaignPreview from '../components/CampaignPreview.vue';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';

export default Vue.extend({
  components: {
    CampaignPreview,
    EmptyPlaceholder,
    MarketingSubnav,
  },

  data() {
    return {
      previewItem: null,
      queryParams: {
        page: 1,
        query: '',
        orderBy: 'created_at',
        order: 'desc',
      },
      pollID: null,
      campaignStatsData: {},
      searchDebounce: null,

      // Table bulk row selection states.
      bulk: {
        checked: [],
        all: false,
      },
    };
  },

  methods: {
    // Campaign statuses.
    canStart(c) {
      return c.status === 'draft' && !c.sendAt;
    },
    canSchedule(c) {
      return c.status === 'draft' && c.sendAt;
    },
    canPause(c) {
      return c.status === 'running';
    },
    canCancel(c) {
      return c.status === 'running' || c.status === 'paused';
    },
    canResume(c) {
      return c.status === 'paused';
    },
    isSheduled(c) {
      return c.status === 'scheduled' || c.sendAt !== null;
    },
    isDone(c) {
      return c.status === 'finished' || c.status === 'cancelled';
    },
    canDownloadReport(c) {
      return ['finished', 'paused', 'running'].includes(c.status) && c.sent > 0;
    },

    isRunning(id) {
      if (id in this.campaignStatsData) {
        return true;
      }
      return false;
    },

    highlightedRow(data) {
      if (data.status === 'running') {
        return ['running'];
      }
      return '';
    },

    recipientsCount(stats) {
      return this.metricVal(stats, ['toSend', 'to_send', 'sent']);
    },

    sentCount(stats) {
      const sent = this.metricVal(stats, ['sent']);
      const toSend = this.metricVal(stats, ['toSend', 'to_send']);
      return Math.max(sent, toSend);
    },

    metricVal(obj, keys) {
      for (let i = 0; i < keys.length; i += 1) {
        const v = obj && obj[keys[i]];
        if (v !== undefined && v !== null && v !== '') {
          const n = Number(v);
          if (!Number.isNaN(n)) {
            return n;
          }
        }
      }
      return 0;
    },

    pctLabel(num, den) {
      if (!den || den <= 0) {
        return '0%';
      }
      const pct = (num / den) * 100;
      if (Number.isInteger(pct)) {
        return `${pct}%`;
      }
      return `${pct.toFixed(2)}%`;
    },

    statusLabel(c) {
      if (c.status === 'finished') return 'Sent';
      if (c.status === 'running') return 'Sending';
      if (c.status === 'scheduled') return 'Scheduled';
      if (c.status === 'paused') return 'Paused';
      if (c.status === 'cancelled') return 'Cancelled';
      return 'Draft';
    },

    statusWhen(c, stats) {
      const raw = (stats && (stats.startedAt || stats.updatedAt))
        || c.startedAt
        || c.sendAt
        || c.updatedAt
        || c.createdAt;
      const d = raw ? dayjs(raw) : null;
      if (!d || !d.isValid()) return '';
      const when = d.format('MMM D, YYYY h:mm A');

      if (c.status === 'finished') {
        return `Sent on ${when}`;
      }
      if (c.status === 'scheduled' && c.sendAt) {
        return `Scheduled for ${dayjs(c.sendAt).format('MMM D, YYYY h:mm A')}`;
      }
      if (c.status === 'running') {
        return when;
      }
      if (c.status === 'paused' || c.status === 'cancelled' || c.status === 'draft') {
        return when;
      }
      return when;
    },

    onSearchInput() {
      clearTimeout(this.searchDebounce);
      this.searchDebounce = setTimeout(() => {
        this.queryParams.page = 1;
        this.getCampaigns();
      }, 250);
    },

    onPageChange(p) {
      if (p < 1 || p > this.totalPages) return;
      this.queryParams.page = p;
      this.getCampaigns();
    },

    onSort(field, direction) {
      this.queryParams.orderBy = field;
      this.queryParams.order = direction;
      this.getCampaigns();
    },

    // Campaign actions.
    previewCampaign(c) {
      this.previewItem = c;
    },

    closePreview() {
      this.previewItem = null;
    },

    getCampaigns() {
      this.$api.getCampaigns({
        page: this.queryParams.page,
        query: this.queryParams.query.replace(/[^\p{L}\p{N}\s]/gu, ' '),
        order_by: this.queryParams.orderBy,
        order: this.queryParams.order,
        no_body: true,
      });
    },

    // Stats returns the campaign object with stats (sent, toSend etc.)
    // if there's live stats available for running campaigns. Otherwise,
    // it returns the incoming campaign object that has the static stats
    // values.
    getCampaignStats(c) {
      if (c.id in this.campaignStatsData) {
        return this.campaignStatsData[c.id];
      }
      return c;
    },

    pollStats() {
      // Clear any running status polls.
      clearInterval(this.pollID);

      // Poll for the status as long as the import is running.
      this.pollID = setInterval(() => {
        this.$api.getCampaignStats().then((data) => {
          // Stop polling. No running campaigns.
          if (data.length === 0) {
            clearInterval(this.pollID);

            // There were running campaigns and stats earlier. Clear them
            // and refetch the campaigns list with up-to-date fields.
            if (Object.keys(this.campaignStatsData).length > 0) {
              this.getCampaigns();
              this.campaignStatsData = {};
            }
          } else {
            // Turn the list of campaigns [{id: 1, ...}, {id: 2, ...}] into
            // a map indexed by the id: {1: {}, 2: {}}.
            this.campaignStatsData = data.reduce((obj, cur) => ({ ...obj, [cur.id]: cur }), {});
          }
        }, () => {
          clearInterval(this.pollID);
        });
      }, 1000);
    },

    changeCampaignStatus(c, status) {
      this.$api.changeCampaignStatus(c.id, status).then(() => {
        this.$utils.toast(this.$t('campaigns.statusChanged', { name: c.name, status }));
        this.getCampaigns();
        this.pollStats();
      });
    },

    async cloneCampaign(name, c) {
      // Fetch the template body from the server.
      let body = '';
      let bodySource = null;
      await this.$api.getCampaign(c.id).then((data) => {
        body = data.body;
        bodySource = data.bodySource;
      });

      const now = this.$utils.getDate();
      const sendLater = !!c.sendAt;
      let sendAt = null;
      if (sendLater) {
        sendAt = dayjs(c.sendAt).isAfter(now) ? c.sendAt : now.add(7, 'day');
      }

      const data = {
        name,
        subject: c.subject,
        lists: c.lists.map((l) => l.id),
        type: c.type,
        from_email: c.fromEmail,
        reply_to: c.replyTo,
        content_type: c.contentType,
        messenger: c.messenger,
        tags: c.tags,
        template_id: c.templateId,
        body,
        body_source: bodySource,
        altbody: c.altbody,
        headers: c.headers,
        send_later: sendLater,
        send_at: sendAt,
        archive: c.archive,
        archive_template_id: c.archiveTemplateId,
        archive_meta: c.archiveMeta,
        media: c.media.map((m) => m.id),
      };

      if (c.archive) {
        data.archive_slug = `${name.toLowerCase().replace(/[^a-z0-9]/g, '-')}-${Date.now().toString().slice(-4)}`;
      }

      this.$api.createCampaign(data).then((d) => {
        this.$router.push({ name: 'campaign', params: { id: d.id } });
      });
    },

    deleteCampaign(c) {
      this.$api.deleteCampaign(c.id).then(() => {
        this.getCampaigns();
        this.$utils.toast(this.$t('globals.messages.deleted', { name: c.name }));
      });
    },

    // Mark all campaigns in the query as selected.
    onSelectAll() {
      this.bulk.all = true;
    },

    onTableCheck() {
      // Disable bulk.all selection if there are no rows checked in the table.
      if (this.bulk.checked.length !== this.campaigns.total) {
        this.bulk.all = false;
      }
    },

    deleteCampaigns() {
      const name = this.$tc('globals.terms.campaign', this.numSelectedCampaigns);

      const fn = () => {
        const params = {};
        if (!this.bulk.all && this.bulk.checked.length > 0) {
          // If 'all' is not selected, delete campaigns by IDs.
          params.id = this.bulk.checked.map((c) => c.id);
        } else {
          // 'All' is selected, delete by query.
          params.query = this.queryParams.query.replace(/[^\p{L}\p{N}\s]/gu, ' ');
          params.all = this.bulk.all;
        }

        this.$api.deleteCampaigns(params)
          .then(() => {
            this.getCampaigns();
            this.$utils.toast(this.$tc(
              'globals.messages.deletedCount',
              this.numSelectedCampaigns,
              { num: this.numSelectedCampaigns, name },
            ));
          });
      };

      this.$utils.confirm(this.$tc(
        'globals.messages.confirmDelete',
        this.numSelectedCampaigns,
        { num: this.numSelectedCampaigns, name: name.toLowerCase() },
      ), fn);
    },
  },

  computed: {
    ...mapState(['campaigns', 'loading']),

    numSelectedCampaigns() {
      return this.bulk.all ? this.campaigns.total : this.bulk.checked.length;
    },

    totalPages() {
      const per = this.campaigns.perPage || 20;
      const total = this.campaigns.total || 0;
      return Math.max(1, Math.ceil(total / per));
    },

    pageRangeLabel() {
      const per = this.campaigns.perPage || 20;
      const total = this.campaigns.total || 0;
      if (!total) return '0 of 0';
      const start = ((this.queryParams.page - 1) * per) + 1;
      const end = Math.min(this.queryParams.page * per, total);
      return `${start}-${end} of ${total}`;
    },
  },

  created() {
    this.$root.$on('page.refresh', this.getCampaigns);
  },

  mounted() {
    this.getCampaigns();
    this.pollStats();
  },

  destroyed() {
    clearInterval(this.pollID);
    clearTimeout(this.searchDebounce);
    this.$root.$off('page.refresh', this.getCampaigns);
  },
});
</script>

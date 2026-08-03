<template>
  <section class="campaigns campaigns-brevo">
    <marketing-subnav />

    <header class="columns page-header">
      <div class="column is-10">
        <h1 class="title is-4">
          {{ $t('globals.terms.campaigns') }}
          <span v-if="!isNaN(campaigns.total)">({{ campaigns.total }})</span>
        </h1>
      </div>
      <div class="column has-text-right">
        <b-field v-if="$can('campaigns:manage')" expanded>
          <b-button expanded :to="{ name: 'campaign', params: { id: 'new' } }" tag="router-link" class="btn-new"
            type="is-primary" icon-left="plus" data-cy="btn-new">
            Create campaign
          </b-button>
        </b-field>
      </div>
    </header>

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
      pagination-position="both"
      @page-change="onPageChange"
      :current-page="queryParams.page"
      :per-page="campaigns.perPage"
      :total="campaigns.total"
      checkable
      backend-sorting
      @sort="onSort"
    >
      <template #top-left>
        <div class="columns">
          <div class="column is-6">
            <form @submit.prevent="getCampaigns">
              <div>
                <b-field>
                  <b-input
                    v-model="queryParams.query"
                    name="query"
                    expanded
                    :placeholder="$t('campaigns.queryPlaceholder')"
                    icon="magnify"
                    ref="query"
                  />
                  <p class="controls">
                    <b-button native-type="submit" type="is-primary" icon-left="magnify" />
                  </p>
                </b-field>
              </div>
            </form>
          </div>
        </div>

        <div class="actions" v-if="bulk.checked.length > 0">
          <a class="a" href="#" @click.prevent="deleteCampaigns" data-cy="btn-delete-campaigns">
            <b-icon icon="trash-can-outline" size="is-small" /> Delete
          </a>
          <span class="a">
            {{ $tc('globals.messages.numSelected', numSelectedCampaigns, { num: numSelectedCampaigns }) }}
            <span v-if="!bulk.all && campaigns.total > campaigns.perPage">
              &mdash;
              <a href="#" @click.prevent="onSelectAll" data-cy="select-all-campaigns">
                {{ $tc('globals.messages.selectAll', campaigns.total, { num: campaigns.total }) }}
              </a>
            </span>
          </span>
        </div>
      </template>

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
              <span>{{ campaignStatusLine(props.row, stats) }}</span>
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
              <strong class="bv-metric__val">{{ $utils.formatNumber(props.row.views || 0) }}</strong>
              <span class="bv-metric__pct">{{ pctLabel(props.row.views || 0, stats.sent || 0) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Clicks</span>
              <strong class="bv-metric__val">{{ $utils.formatNumber(props.row.clicks || 0) }}</strong>
              <span class="bv-metric__pct">{{ pctLabel(props.row.clicks || 0, stats.sent || 0) }}</span>
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
            <b-dropdown class="campaign-actions-menu" position="is-top-left">
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

                <router-link
                  class="campaign-action"
                  :to="{ name: 'campaign', params: { id: props.row.id }, query: { view: 'report' } }"
                  data-cy="btn-campaign-view-report"
                  aria-label="Report"
                >
                  <b-tooltip label="Report" type="is-dark" position="is-left">
                    <b-icon icon="chart-bar" />
                  </b-tooltip>
                </router-link>

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
      return stats.toSend || stats.sent || 0;
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

    campaignStatusLine(c, stats) {
      const raw = (stats && (stats.startedAt || stats.updatedAt))
        || c.startedAt
        || c.sendAt
        || c.updatedAt
        || c.createdAt;
      const d = raw ? dayjs(raw) : null;
      const when = d ? d.format('MMM D, YYYY h:mm A') : '';

      if (c.status === 'finished') {
        return when ? `Sent on ${when}` : 'Sent';
      }
      if (c.status === 'running') {
        return when ? `Sending · ${when}` : 'Sending';
      }
      if (c.status === 'scheduled' && c.sendAt) {
        return `Scheduled for ${dayjs(c.sendAt).format('MMM D, YYYY h:mm A')}`;
      }
      if (c.status === 'paused') {
        return when ? `Paused · ${when}` : 'Paused';
      }
      if (c.status === 'cancelled') {
        return when ? `Cancelled · ${when}` : 'Cancelled';
      }
      return when ? `Draft · ${when}` : 'Draft';
    },

    onPageChange(p) {
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
    this.$root.$off('page.refresh', this.getCampaigns);
  },
});
</script>

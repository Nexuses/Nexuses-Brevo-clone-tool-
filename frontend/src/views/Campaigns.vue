<template>
  <section class="campaigns campaigns-brevo">
    <marketing-subnav />

    <header class="campaigns-brevo__header">
      <div class="campaigns-brevo__title-row">
        <h1 class="campaigns-brevo__title">Campaigns</h1>
        <div class="campaigns-brevo__actions">
          <button type="button" class="campaigns-brevo__btn campaigns-brevo__btn--outline" @click="onCreateFolder">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path d="M2.5 4.5h4l1.2 1.5H13.5v6.5a1 1 0 01-1 1h-10a1 1 0 01-1-1V5.5a1 1 0 011-1z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round" />
            </svg>
            Create folder
          </button>
          <button
            v-if="$can('campaigns:manage')"
            type="button"
            class="campaigns-brevo__btn campaigns-brevo__btn--primary"
            data-cy="btn-new"
            @click="openCreateModal('standard')"
          >
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path d="M6 2v8M2 6h8" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
            Create campaign
          </button>
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
        <label class="campaigns-brevo__check">
          <input
            type="checkbox"
            :checked="allVisibleChecked"
            :indeterminate.prop="someVisibleChecked"
            aria-label="Select all campaigns"
            @change="toggleSelectAllVisible"
          />
        </label>
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

        <div class="campaigns-brevo__dd" ref="statusDd">
          <button type="button" class="campaigns-brevo__chip" @click.stop="toggleStatusMenu">
            {{ statusFilterLabel }}
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
          <ul v-if="statusMenuOpen" class="campaigns-brevo__menu">
            <li v-for="opt in statusOptions" :key="opt.value || 'all'">
              <button type="button" class="campaigns-brevo__menu-item" :class="{ 'is-active': queryParams.status === opt.value }" @click="setStatusFilter(opt.value)">
                <span v-if="opt.dot" class="campaigns-brevo__menu-dot" :class="'is-' + opt.dot" />
                {{ opt.label }}
              </button>
            </li>
          </ul>
        </div>

        <div class="campaigns-brevo__dd" ref="tagDd">
          <button type="button" class="campaigns-brevo__chip" @click.stop="toggleTagMenu">
            {{ tagFilterLabel }}
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
          <ul v-if="tagMenuOpen" class="campaigns-brevo__menu">
            <li>
              <button type="button" class="campaigns-brevo__menu-item" :class="{ 'is-active': !queryParams.tag }" @click="setTagFilter(null)">
                All tags
              </button>
            </li>
            <li v-for="tag in availableTags" :key="tag">
              <button type="button" class="campaigns-brevo__menu-item" :class="{ 'is-active': queryParams.tag === tag }" @click="setTagFilter(tag)">
                {{ tag }}
              </button>
            </li>
            <li v-if="!availableTags.length" class="campaigns-brevo__menu-empty">No tags yet</li>
          </ul>
        </div>
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
        <label class="campaigns-brevo__pager-jump">
          <select :value="queryParams.page" aria-label="Page" @change="onPageSelect">
            <option v-for="p in totalPages" :key="p" :value="p">{{ p }}</option>
          </select>
        </label>
        <span class="campaigns-brevo__pager-pages">of {{ totalPages }} pages</span>
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
        <div
          class="bv-campaign-card"
          :class="{ 'is-clickable': true }"
          :set="stats = getCampaignStats(props.row)"
          role="link"
          tabindex="0"
          @click="openCampaign(props.row, $event)"
          @keydown.enter.prevent="openCampaign(props.row, $event)"
        >
          <div class="bv-campaign-card__main">
            <router-link
              class="bv-campaign-card__title"
              :to="campaignTarget(props.row)"
              @click.native.stop
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
              <strong class="bv-metric__val">{{ metricOrDash(props.row, recipientsCount(stats)) }}</strong>
              <span class="bv-metric__pct">{{ metricPctOrDash(props.row, recipientsCount(stats), recipientsCount(stats)) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Opens</span>
              <strong class="bv-metric__val">{{ metricOrDash(props.row, metricVal(stats, ['views', 'opens', 'uniqueOpens', 'unique_opens'])) }}</strong>
              <span class="bv-metric__pct">{{ metricPctOrDash(props.row, metricVal(stats, ['views', 'opens', 'uniqueOpens', 'unique_opens']), sentCount(stats)) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Clicks</span>
              <strong class="bv-metric__val">{{ metricOrDash(props.row, metricVal(stats, ['clicks'])) }}</strong>
              <span class="bv-metric__pct">{{ metricPctOrDash(props.row, metricVal(stats, ['clicks']), sentCount(stats)) }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Unsubscribed</span>
              <strong class="bv-metric__val">{{ hasMetrics(props.row) ? '0' : '—' }}</strong>
              <span class="bv-metric__pct">{{ hasMetrics(props.row) ? '0%' : '' }}</span>
            </div>
            <div class="bv-metric">
              <span class="bv-metric__lbl">Conversions</span>
              <strong class="bv-metric__val">—</strong>
              <span class="bv-metric__pct">{{ hasMetrics(props.row) ? '0%' : '' }}</span>
            </div>
          </div>

          <div class="bv-campaign-card__actions" @click.stop>
            <router-link
              v-if="isEditableCampaign(props.row)"
              class="bv-campaign-card__edit"
              :to="campaignTarget(props.row)"
              data-cy="btn-campaign-edit"
              aria-label="Edit"
            >
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
                <path d="M3.5 12.8l8.4-8.4 2.7 2.7-8.4 8.4H3.5v-2.7z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" />
                <path d="M10.6 5.7l1.7 1.7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            </router-link>
            <router-link
              v-else
              class="bv-campaign-card__report"
              :to="campaignTarget(props.row)"
              data-cy="btn-campaign-view-report"
              aria-label="Report"
            >
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
                <path d="M3 13.5l3.2-3.2 2.3 2.3L15 6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
                <path d="M12 6h3v3" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
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

    <b-modal
      :active.sync="isCreateOpen"
      :width="760"
      scroll="keep"
      class="create-campaign-form-modal"
    >
      <create-campaign-modal
        v-if="isCreateOpen"
        :start-section="createModalSection"
        @close="isCreateOpen = false"
      />
    </b-modal>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';
import CampaignPreview from '../components/CampaignPreview.vue';
import CreateCampaignModal from '../components/CreateCampaignModal.vue';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';

export default Vue.extend({
  components: {
    CampaignPreview,
    CreateCampaignModal,
    EmptyPlaceholder,
    MarketingSubnav,
  },

  data() {
    return {
      previewItem: null,
      isCreateOpen: false,
      createModalSection: 'standard',
      queryParams: {
        page: 1,
        query: '',
        status: '',
        tag: null,
        orderBy: 'created_at',
        order: 'desc',
      },
      pollID: null,
      campaignStatsData: {},
      searchDebounce: null,
      statusMenuOpen: false,
      tagMenuOpen: false,
      knownTags: [],
      statusOptions: [
        { value: '', label: 'All statuses' },
        { value: 'draft', label: 'Draft', dot: 'draft' },
        { value: 'finished', label: 'Sent', dot: 'sent' },
        { value: 'scheduled', label: 'Scheduled', dot: 'scheduled' },
        { value: 'running', label: 'Running', dot: 'running' },
        { value: 'paused', label: 'Suspended', dot: 'suspended' },
        { value: 'cancelled', label: 'Archived', dot: 'archived' },
      ],

      // Table bulk row selection states.
      bulk: {
        checked: [],
        all: false,
      },
    };
  },

  methods: {
    openCreateModal(section) {
      this.createModalSection = section || 'standard';
      this.isCreateOpen = true;
    },

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
      if (c.status === 'running') return 'Running';
      if (c.status === 'scheduled') return 'Scheduled';
      if (c.status === 'paused') return 'Suspended';
      if (c.status === 'cancelled') return 'Archived';
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

      if (c.status === 'finished') return `Sent on ${when}`;
      if (c.status === 'scheduled' && c.sendAt) {
        return `Scheduled for ${dayjs(c.sendAt).format('MMM D, YYYY h:mm A')}`;
      }
      if (c.status === 'running') return when;
      if (c.status === 'draft') return `Last edited ${when}`;
      if (c.status === 'paused' || c.status === 'cancelled') return `Last edited ${when}`;
      return when;
    },

    hasMetrics(c) {
      return ['finished', 'running', 'paused', 'cancelled'].indexOf(c.status) > -1;
    },

    metricOrDash(c, val) {
      if (!this.hasMetrics(c)) return '—';
      return this.$utils.formatNumber(val || 0);
    },

    metricPctOrDash(c, num, den) {
      if (!this.hasMetrics(c)) return '';
      return this.pctLabel(num, den);
    },

    isEditableCampaign(c) {
      return c.status === 'draft' || c.status === 'scheduled' || c.status === 'paused';
    },

    isReportCampaign(c) {
      return c.status === 'finished' || c.status === 'cancelled' || c.status === 'running';
    },

    campaignTarget(c) {
      if (this.isReportCampaign(c)) {
        return { name: 'campaign', params: { id: c.id }, query: { view: 'report' } };
      }
      return { name: 'campaign', params: { id: c.id } };
    },

    openCampaign(c, e) {
      if (e && e.target && e.target.closest) {
        if (e.target.closest('.bv-campaign-card__actions, .checkbox-cell, input, button, .dropdown, a')) {
          return;
        }
      }
      this.$router.push(this.campaignTarget(c));
    },

    onCreateFolder() {
      this.$utils.toast('Campaign folders coming soon');
    },

    toggleStatusMenu() {
      this.tagMenuOpen = false;
      this.statusMenuOpen = !this.statusMenuOpen;
    },

    toggleTagMenu() {
      this.statusMenuOpen = false;
      this.tagMenuOpen = !this.tagMenuOpen;
    },

    setStatusFilter(status) {
      this.queryParams.status = status || '';
      this.queryParams.page = 1;
      this.statusMenuOpen = false;
      this.getCampaigns();
    },

    setTagFilter(tag) {
      this.queryParams.tag = tag;
      this.queryParams.page = 1;
      this.tagMenuOpen = false;
      this.getCampaigns();
    },

    toggleSelectAllVisible(e) {
      if (e.target.checked) {
        this.bulk.checked = [...(this.campaigns.results || [])];
      } else {
        this.bulk.checked = [];
        this.bulk.all = false;
      }
    },

    onPageSelect(e) {
      this.onPageChange(parseInt(e.target.value, 10) || 1);
    },

    onDocClick(e) {
      if (this.statusMenuOpen && this.$refs.statusDd && !this.$refs.statusDd.contains(e.target)) {
        this.statusMenuOpen = false;
      }
      if (this.tagMenuOpen && this.$refs.tagDd && !this.$refs.tagDd.contains(e.target)) {
        this.tagMenuOpen = false;
      }
    },

    rememberTags(rows) {
      const next = new Set(this.knownTags);
      (rows || []).forEach((c) => {
        (c.tags || []).forEach((t) => {
          if (t) next.add(t);
        });
      });
      this.knownTags = Array.from(next).sort();
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
      const params = {
        page: this.queryParams.page,
        query: this.queryParams.query.replace(/[^\p{L}\p{N}\s]/gu, ' '),
        order_by: this.queryParams.orderBy,
        order: this.queryParams.order,
        no_body: true,
      };
      if (this.queryParams.status) {
        params.status = this.queryParams.status;
      }
      if (this.queryParams.tag) {
        params.tag = this.queryParams.tag;
      }
      this.$api.getCampaigns(params).then((data) => {
        this.rememberTags((data && data.results) || this.campaigns.results || []);
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

    allVisibleChecked() {
      const rows = this.campaigns.results || [];
      return rows.length > 0 && this.bulk.checked.length === rows.length;
    },

    someVisibleChecked() {
      const rows = this.campaigns.results || [];
      return this.bulk.checked.length > 0 && this.bulk.checked.length < rows.length;
    },

    statusFilterLabel() {
      const opt = this.statusOptions.find((o) => o.value === this.queryParams.status);
      return (opt && opt.label) || 'All statuses';
    },

    tagFilterLabel() {
      return this.queryParams.tag || 'Select tags';
    },

    availableTags() {
      return this.knownTags;
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
    document.addEventListener('click', this.onDocClick);
    this.getCampaigns();
    this.pollStats();
    if (this.$route.query.create) {
      const section = this.$route.query.section === 'automated' ? 'automated' : 'standard';
      this.openCreateModal(section);
    }
  },

  destroyed() {
    document.removeEventListener('click', this.onDocClick);
    clearInterval(this.pollID);
    clearTimeout(this.searchDebounce);
    this.$root.$off('page.refresh', this.getCampaigns);
  },
});
</script>

<template>
  <section class="campaign-report" data-ui="brevo-report">
    <header class="campaign-report__top">
      <router-link :to="{ name: 'campaigns' }" class="campaign-report__back" aria-label="Back to campaigns">
        <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
          <path
            d="M11.5 4.5L7 9l4.5 4.5"
            stroke="currentColor"
            stroke-width="1.7"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        <span>Back</span>
      </router-link>

      <div class="campaign-report__thumb" aria-hidden="true">
        <iframe
          v-if="campaign.id"
          class="campaign-report__thumb-frame"
          :src="`/api/campaigns/${campaign.id}/preview`"
          title="Campaign preview"
          tabindex="-1"
        />
        <div v-else class="campaign-report__thumb-inner">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="1.5" />
            <path d="M3.5 7l8.5 6 8.5-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          </svg>
        </div>
      </div>

      <div class="campaign-report__heading">
        <h1 class="campaign-report__title">{{ campaign.name }}</h1>
        <p class="campaign-report__sub">
          #{{ campaign.id }}
          <span v-if="sentLabel"> · {{ sentLabel }}</span>
        </p>
        <div class="campaign-report__meta">
          <div>
            <span class="lbl">Subject</span>
            <strong>{{ campaign.subject || '—' }}</strong>
          </div>
          <div>
            <span class="lbl">From</span>
            <strong>{{ campaign.fromEmail || '—' }}</strong>
          </div>
          <div>
            <span class="lbl">Reply to</span>
            <strong>{{ campaign.replyTo || campaign.fromEmail || '—' }}</strong>
          </div>
        </div>
      </div>

      <div class="campaign-report__actions">
        <button
          type="button"
          class="campaign-report__icon-btn"
          aria-label="Copy report link"
          title="Copy report link"
          @click="onShare"
        >
          <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
            <rect
              x="6.25"
              y="6.25"
              width="8.5"
              height="8.5"
              rx="1.5"
              stroke="currentColor"
              stroke-width="1.5"
            />
            <path
              d="M4.5 11.75H3.75A1.5 1.5 0 012.25 10.25V3.75A1.5 1.5 0 013.75 2.25h6.5A1.5 1.5 0 0111.75 3.75V4.5"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
        </button>
        <a
          v-if="canExport"
          class="campaign-report__export-btn"
          :href="`/api/campaigns/${campaign.id}/report`"
          target="_blank"
          rel="noopener noreferrer"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
            <path
              d="M8 2.5v7M5 7l3 3 3-3M3 12.5h10"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          Export report
        </a>
      </div>
    </header>

    <nav class="campaign-report__tabs" aria-label="Report sections">
      <button
        v-for="t in tabs"
        :key="t.id"
        type="button"
        class="campaign-report__tab"
        :class="{ 'is-active': activeTab === t.id }"
        @click="activeTab = t.id"
      >
        {{ t.label }}
      </button>
    </nav>

    <!-- Overview -->
    <div v-if="activeTab === 'overview'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Campaign performance</h2>
        <span class="hint">
          Automated opens and clicks included
          <b-tooltip label="Includes automated and privacy-proxy opens and clicks.">
            <b-icon icon="help-circle-outline" size="is-small" />
          </b-tooltip>
        </span>
      </div>

      <div class="campaign-report__metrics">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Delivered</span>
            <router-link
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              :to="subscribersRoute"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(delivered) }}</div>
          <div class="metric__rate">Delivery rate {{ deliveryRate }}</div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Opens
              <b-tooltip label="Unique contacts who opened this campaign.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
            <button type="button" class="metric__view" @click="activeTab = 'opens'">
              <b-icon icon="account-outline" size="is-small" /> View
            </button>
          </div>
          <div class="metric__value">{{ n(uniqueOpens) }}</div>
          <div class="metric__rate">Open rate {{ openRate }}</div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Clicks
              <b-tooltip label="Unique contacts who clicked a link.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
            <button type="button" class="metric__view" @click="activeTab = 'clicks'">
              <b-icon icon="account-outline" size="is-small" /> View
            </button>
          </div>
          <div class="metric__value">{{ n(uniqueClicks) }}</div>
          <div class="metric__rate">Click-through rate {{ clickRate }}</div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Conversions</span>
          </div>
          <div class="metric__value">0</div>
          <div class="metric__rate">Conversion rate 0%</div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Unsubscribes</span>
            <button type="button" class="metric__view" @click="activeTab = 'unsubscribes'">
              <b-icon icon="account-outline" size="is-small" /> View
            </button>
          </div>
          <div class="metric__value">{{ n(unsubscribes) }}</div>
          <div class="metric__rate">Unsubscribe rate {{ unsubRate }}</div>
        </div>
      </div>

      <h2 class="campaign-report__h2">Campaign audience</h2>
      <div class="campaign-report__audience">
        <button
          type="button"
          class="audience__head"
          :aria-expanded="audienceOpen ? 'true' : 'false'"
          @click="audienceOpen = !audienceOpen"
        >
          <b-icon :icon="audienceOpen ? 'minus' : 'plus'" size="is-small" />
          Included lists
        </button>
        <template v-if="audienceOpen">
          <div v-if="audienceLists.length === 0" class="audience__empty">
            No lists attached to this campaign.
          </div>
          <div
            v-for="l in audienceLists"
            :key="l.id"
            class="audience__row"
          >
            <div class="audience__name">{{ l.name }}</div>
            <div class="audience__count">{{ n(l.subscriberCount || 0) }} contacts</div>
            <router-link
              :to="{ name: 'subscribers_list', params: { listID: l.id } }"
              class="audience__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
        </template>
      </div>
    </div>

    <!-- Deliverability -->
    <div v-else-if="activeTab === 'deliverability'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Deliverability details</h2>
      </div>

      <div class="campaign-report__metrics campaign-report__metrics--6">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Sent to</span>
          </div>
          <div class="metric__value">{{ n(sentTo) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Delivered</span>
            <router-link
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              :to="subscribersRoute"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(delivered) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Delivery rate</span>
          </div>
          <div class="metric__value">{{ deliveryRate }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              In Processing
              <b-tooltip label="Emails still queued or sending.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
          </div>
          <div class="metric__value">{{ n(inProcessing) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Soft bounces
              <b-tooltip label="Temporary delivery failures such as a full inbox.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
            <router-link
              :to="{ name: 'bounces', query: { campaign_id: campaign.id, type: 'soft' } }"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(softBounces) }}</div>
          <div class="metric__rate">{{ softRate }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Hard bounces
              <b-tooltip label="Permanent delivery failures such as an invalid address.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
            <router-link
              :to="{ name: 'bounces', query: { campaign_id: campaign.id, type: 'hard' } }"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(hardBounces) }}</div>
          <div class="metric__rate">{{ hardRate }}</div>
        </div>
      </div>

      <div class="campaign-report__section-head campaign-report__section-head--table">
        <h2 class="campaign-report__inline-title">
          Reasons for
          <select v-model="bounceReasonType" class="campaign-report__select" aria-label="Bounce type">
            <option value="soft">soft bounce</option>
            <option value="hard">hard bounce</option>
            <option value="complaint">complaint</option>
          </select>
          <b-tooltip :label="bounceReasonHelp">
            <b-icon icon="help-circle-outline" size="is-small" />
          </b-tooltip>
        </h2>
        <a
          v-if="canExport"
          :href="`/api/campaigns/${campaign.id}/report`"
          class="campaign-report__export"
          target="_blank"
          rel="noopener noreferrer"
        >
          <b-icon icon="download-outline" size="is-small" /> Export
        </a>
      </div>

      <div class="campaign-report__table-wrap">
        <table class="campaign-report__table">
          <thead>
            <tr>
              <th>Reason</th>
              <th class="is-num">{{ bounceReasonLabel }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="bounceReasonRows.length === 0">
              <td colspan="2" class="is-empty">No {{ bounceReasonLabel.toLowerCase() }} for this campaign.</td>
            </tr>
            <tr v-for="(row, i) in bounceReasonRows" :key="`${row.reason}-${i}`">
              <td>{{ row.reason }}</td>
              <td class="is-num">
                <strong>{{ n(row.count) }}</strong>
                <span class="sub">{{ pct(row.count, sentTo) }} of recipients</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="campaign-report__section-head campaign-report__section-head--table">
        <h2 class="campaign-report__inline-title">
          Deliverability breakdown by
          <select v-model="breakdownBy" class="campaign-report__select" aria-label="Breakdown">
            <option value="lists">lists</option>
          </select>
        </h2>
        <a
          v-if="canExport"
          :href="`/api/campaigns/${campaign.id}/report`"
          class="campaign-report__export"
          target="_blank"
          rel="noopener noreferrer"
        >
          <b-icon icon="download-outline" size="is-small" /> Export
        </a>
      </div>

      <div class="campaign-report__table-wrap campaign-report__table-wrap--scroll">
        <table class="campaign-report__table campaign-report__table--wide">
          <thead>
            <tr>
              <th>List name</th>
              <th class="is-num">Delivery rate</th>
              <th class="is-num">Processing</th>
              <th class="is-num">Deferred</th>
              <th class="is-num">Soft bounces</th>
              <th class="is-num">Hard bounces</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="audienceLists.length === 0">
              <td colspan="6" class="is-empty">No lists attached to this campaign.</td>
            </tr>
            <tr v-for="l in audienceLists" :key="`del-${l.id}`">
              <td>
                <div class="cell-name">{{ l.name }}</div>
                <div class="cell-meta">
                  #{{ l.id }} · {{ n(l.subscriberCount || 0) }} contacts - {{ listShare(l) }} of all
                </div>
              </td>
              <td class="is-num">
                <strong>{{ deliveryRate }}</strong>
                <span class="sub">{{ n(listStat(l, delivered)) }} contacts</span>
              </td>
              <td class="is-num">
                <strong>{{ processingRate }}</strong>
                <span class="sub">{{ n(listStat(l, inProcessing)) }} contact{{ inProcessing === 1 ? '' : 's' }}</span>
              </td>
              <td class="is-num">
                <strong>0%</strong>
                <span class="sub">0 contacts</span>
              </td>
              <td class="is-num">
                <strong>{{ softRate }}</strong>
                <span class="sub">{{ n(listStat(l, softBounces)) }} contacts</span>
              </td>
              <td class="is-num">
                <strong>{{ hardRate }}</strong>
                <span class="sub">{{ n(listStat(l, hardBounces)) }} contacts</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Opens -->
    <div v-else-if="activeTab === 'opens'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Opens Details</h2>
        <span class="hint">
          Automated opens included
          <b-tooltip label="Includes automated and privacy-proxy opens.">
            <b-icon icon="help-circle-outline" size="is-small" />
          </b-tooltip>
        </span>
      </div>

      <div class="campaign-report__metrics campaign-report__metrics--4">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Opens</span>
            <router-link
              :to="{ name: 'campaignAnalytics', query: { id: campaign.id } }"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(uniqueOpens) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Open rate</span>
          </div>
          <div class="metric__value">{{ openRate }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Total opens</span>
          </div>
          <div class="metric__value">{{ n(totalOpens) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Apple MPP opens
              <b-tooltip label="Opens attributed to Apple Mail Privacy Protection are not tracked separately.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
          </div>
          <div class="metric__value">0</div>
        </div>
      </div>

      <div class="campaign-report__section-head campaign-report__section-head--table">
        <h2 class="campaign-report__inline-title">
          Opens breakdown by
          <select v-model="breakdownBy" class="campaign-report__select" aria-label="Breakdown">
            <option value="lists">lists</option>
          </select>
        </h2>
        <div class="campaign-report__table-actions">
          <span class="hint">Bot opens excluded.</span>
          <a
            v-if="canExport"
            :href="`/api/campaigns/${campaign.id}/report`"
            class="campaign-report__export"
            target="_blank"
            rel="noopener noreferrer"
          >
            <b-icon icon="download-outline" size="is-small" /> Export
          </a>
        </div>
      </div>

      <div class="campaign-report__table-wrap">
        <table class="campaign-report__table">
          <thead>
            <tr>
              <th>List name</th>
              <th class="is-num">Open rate</th>
              <th class="is-num">Total opens</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="audienceLists.length === 0">
              <td colspan="3" class="is-empty">No lists attached to this campaign.</td>
            </tr>
            <tr v-for="l in audienceLists" :key="`open-${l.id}`">
              <td>
                <div class="cell-name">{{ l.name }}</div>
                <div class="cell-meta">
                  #{{ l.id }} · {{ n(l.subscriberCount || 0) }} contacts - {{ listShare(l) }} of all
                </div>
              </td>
              <td class="is-num">
                <strong>{{ openRate }}</strong>
                <span class="sub">{{ n(listStat(l, uniqueOpens)) }} opens</span>
              </td>
              <td class="is-num">
                <strong>{{ totalOpenRate }}</strong>
                <span class="sub">{{ n(listStat(l, totalOpens)) }} opens</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Clicks -->
    <div v-else-if="activeTab === 'clicks'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Clicks details</h2>
        <span class="hint">
          Bot clicks included
          <b-tooltip label="Automated and bot clicks are included in campaign totals.">
            <b-icon icon="help-circle-outline" size="is-small" />
          </b-tooltip>
        </span>
      </div>

      <div class="campaign-report__metrics campaign-report__metrics--4">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Click-through rate</span>
          </div>
          <div class="metric__value">{{ clickRate }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Total clicks
              <b-tooltip label="All click events, including repeats from the same contact.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
          </div>
          <div class="metric__value">{{ n(totalClicks) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Clicks</span>
            <router-link
              :to="{ name: 'campaignAnalytics', query: { id: campaign.id } }"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(uniqueClicks) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">
              Click-to-open rate
              <b-tooltip label="Unique clicks divided by unique opens.">
                <b-icon icon="help-circle-outline" size="is-small" />
              </b-tooltip>
            </span>
          </div>
          <div class="metric__value">{{ ctorRate }}</div>
        </div>
      </div>

      <div class="campaign-report__section-head campaign-report__section-head--table">
        <h2 class="campaign-report__inline-title">
          Clicks breakdown by
          <select v-model="breakdownBy" class="campaign-report__select" aria-label="Breakdown">
            <option value="lists">lists</option>
          </select>
        </h2>
        <div class="campaign-report__table-actions">
          <span class="hint">Bot clicks are excluded</span>
          <a
            v-if="canExport"
            :href="`/api/campaigns/${campaign.id}/report`"
            class="campaign-report__export"
            target="_blank"
            rel="noopener noreferrer"
          >
            <b-icon icon="download-outline" size="is-small" /> Export
          </a>
        </div>
      </div>

      <div class="campaign-report__table-wrap">
        <table class="campaign-report__table">
          <thead>
            <tr>
              <th>List name</th>
              <th class="is-num">Clicks percentage</th>
              <th class="is-num">Total clicks</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="audienceLists.length === 0">
              <td colspan="3" class="is-empty">No lists attached to this campaign.</td>
            </tr>
            <tr v-for="l in audienceLists" :key="`click-${l.id}`">
              <td>
                <div class="cell-name">{{ l.name }}</div>
                <div class="cell-meta">
                  #{{ l.id }} · {{ n(l.subscriberCount || 0) }} contacts - {{ listShare(l) }} of all
                </div>
              </td>
              <td class="is-num">
                <strong>{{ clickRate }}</strong>
                <span class="sub">{{ n(listStat(l, uniqueClicks)) }} clicks</span>
              </td>
              <td class="is-num">
                <strong>{{ totalClickRate }}</strong>
                <span class="sub">{{ n(listStat(l, totalClicks)) }} clicks</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Conversions -->
    <div v-else-if="activeTab === 'conversions'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Conversions details</h2>
      </div>

      <div class="campaign-report__metrics campaign-report__metrics--4">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Conversions</span>
          </div>
          <div class="metric__value">0</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Conversion rate</span>
          </div>
          <div class="metric__value">0%</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Revenue</span>
          </div>
          <div class="metric__value">0</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Orders</span>
          </div>
          <div class="metric__value">0</div>
        </div>
      </div>

      <div class="campaign-report__section-head campaign-report__section-head--table">
        <h2 class="campaign-report__inline-title">
          Conversions breakdown by
          <select v-model="breakdownBy" class="campaign-report__select" aria-label="Breakdown">
            <option value="lists">lists</option>
          </select>
        </h2>
      </div>

      <div class="campaign-report__table-wrap">
        <table class="campaign-report__table">
          <thead>
            <tr>
              <th>List name</th>
              <th class="is-num">Conversions</th>
              <th class="is-num">Conversion rate</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="audienceLists.length === 0">
              <td colspan="3" class="is-empty">No lists attached to this campaign.</td>
            </tr>
            <tr v-for="l in audienceLists" :key="`conv-${l.id}`">
              <td>
                <div class="cell-name">{{ l.name }}</div>
                <div class="cell-meta">
                  #{{ l.id }} · {{ n(l.subscriberCount || 0) }} contacts - {{ listShare(l) }} of all
                </div>
              </td>
              <td class="is-num">
                <strong>0</strong>
                <span class="sub">0 conversions</span>
              </td>
              <td class="is-num">
                <strong>0%</strong>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Unsubscribes -->
    <div v-else-if="activeTab === 'unsubscribes'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Unsubscribes details</h2>
      </div>

      <div class="campaign-report__metrics campaign-report__metrics--4">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Unsubscribes</span>
            <router-link
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              :to="subscribersRoute"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(unsubscribes) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Unsubscribe rate</span>
          </div>
          <div class="metric__value">{{ unsubRate }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Complaints</span>
            <router-link
              :to="{ name: 'bounces', query: { campaign_id: campaign.id, type: 'complaint' } }"
              class="metric__view"
            >
              <b-icon icon="account-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ n(complaints) }}</div>
        </div>
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Complaint rate</span>
          </div>
          <div class="metric__value">{{ complaintRate }}</div>
        </div>
      </div>

      <div class="campaign-report__section-head campaign-report__section-head--table">
        <h2 class="campaign-report__inline-title">
          Unsubscribes breakdown by
          <select v-model="breakdownBy" class="campaign-report__select" aria-label="Breakdown">
            <option value="lists">lists</option>
          </select>
        </h2>
        <a
          v-if="canExport"
          :href="`/api/campaigns/${campaign.id}/report`"
          class="campaign-report__export"
          target="_blank"
          rel="noopener noreferrer"
        >
          <b-icon icon="download-outline" size="is-small" /> Export
        </a>
      </div>

      <div class="campaign-report__table-wrap">
        <table class="campaign-report__table">
          <thead>
            <tr>
              <th>List name</th>
              <th class="is-num">Unsubscribes</th>
              <th class="is-num">Unsubscribe rate</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="audienceLists.length === 0">
              <td colspan="3" class="is-empty">No lists attached to this campaign.</td>
            </tr>
            <tr v-for="l in audienceLists" :key="`unsub-${l.id}`">
              <td>
                <div class="cell-name">{{ l.name }}</div>
                <div class="cell-meta">
                  #{{ l.id }} · {{ n(l.subscriberCount || 0) }} contacts - {{ listShare(l) }} of all
                </div>
              </td>
              <td class="is-num">
                <strong>{{ n(listStat(l, unsubscribes)) }}</strong>
                <span class="sub">{{ n(listStat(l, unsubscribes)) }} contacts</span>
              </td>
              <td class="is-num">
                <strong>{{ unsubRate }}</strong>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';

export default Vue.extend({
  name: 'CampaignReport',

  props: {
    campaign: { type: Object, required: true },
  },

  data() {
    return {
      activeTab: 'overview',
      audienceOpen: true,
      breakdownBy: 'lists',
      bounceReasonType: 'soft',
      summary: {
        uniqueOpens: 0,
        totalOpens: 0,
        uniqueClicks: 0,
        totalClicks: 0,
        unsubscribes: 0,
        softBounces: 0,
        hardBounces: 0,
        complaints: 0,
        bounces: 0,
      },
      bounceReasons: [],
      tabs: [
        { id: 'overview', label: 'Overview' },
        { id: 'deliverability', label: 'Deliverability' },
        { id: 'opens', label: 'Opens' },
        { id: 'clicks', label: 'Clicks' },
        { id: 'conversions', label: 'Conversions' },
        { id: 'unsubscribes', label: 'Unsubscribes' },
      ],
    };
  },

  computed: {
    ...mapState(['lists']),

    canExport() {
      return ['finished', 'paused', 'running'].includes(this.campaign.status)
        && (this.campaign.sent || 0) > 0;
    },

    sentTo() {
      const sent = Number(this.campaign.sent) || 0;
      const toSend = Number(this.campaign.toSend) || 0;
      return Math.max(sent, toSend);
    },

    inProcessing() {
      const sent = Number(this.campaign.sent) || 0;
      const toSend = Number(this.campaign.toSend) || 0;
      return Math.max(0, toSend - sent);
    },

    softBounces() {
      return Number(this.summary.softBounces) || 0;
    },

    hardBounces() {
      return Number(this.summary.hardBounces) || 0;
    },

    complaints() {
      return Number(this.summary.complaints) || 0;
    },

    delivered() {
      return Math.max(0, this.sentTo - this.softBounces - this.hardBounces);
    },

    uniqueOpens() {
      return Number(this.summary.uniqueOpens) || Number(this.campaign.views) || 0;
    },

    totalOpens() {
      return Number(this.summary.totalOpens) || Number(this.campaign.views) || 0;
    },

    uniqueClicks() {
      return Number(this.summary.uniqueClicks) || Number(this.campaign.clicks) || 0;
    },

    totalClicks() {
      return Number(this.summary.totalClicks) || Number(this.campaign.clicks) || 0;
    },

    unsubscribes() {
      return Number(this.summary.unsubscribes) || 0;
    },

    deliveryRate() {
      return this.pct(this.delivered, this.sentTo);
    },

    processingRate() {
      return this.pct(this.inProcessing, this.sentTo);
    },

    softRate() {
      return this.pct(this.softBounces, this.sentTo);
    },

    hardRate() {
      return this.pct(this.hardBounces, this.sentTo);
    },

    openRate() {
      return this.pct(this.uniqueOpens, this.delivered);
    },

    totalOpenRate() {
      return this.pct(this.totalOpens, this.delivered);
    },

    clickRate() {
      return this.pct(this.uniqueClicks, this.delivered);
    },

    totalClickRate() {
      return this.pct(this.totalClicks, this.delivered);
    },

    ctorRate() {
      return this.pct(this.uniqueClicks, this.uniqueOpens);
    },

    unsubRate() {
      return this.pct(this.unsubscribes, this.delivered);
    },

    complaintRate() {
      return this.pct(this.complaints, this.sentTo);
    },

    bounceReasonLabel() {
      if (this.bounceReasonType === 'hard') return 'Hard bounces';
      if (this.bounceReasonType === 'complaint') return 'Complaints';
      return 'Soft bounces';
    },

    bounceReasonHelp() {
      if (this.bounceReasonType === 'hard') {
        return 'Permanent failures such as unknown users or blocked domains.';
      }
      if (this.bounceReasonType === 'complaint') {
        return 'Recipients who marked the campaign as spam.';
      }
      return 'Temporary failures such as a full mailbox or a timeout.';
    },

    bounceReasonRows() {
      return (this.bounceReasons || []).filter((r) => r.type === this.bounceReasonType);
    },

    sentLabel() {
      const c = this.campaign;
      const raw = c.startedAt || c.sendAt || c.updatedAt;
      if (!raw) return '';
      const d = dayjs(raw);
      if (c.status === 'finished' || c.status === 'running') {
        return `Sent on ${d.format('MMM D, YYYY HH:mm')}`;
      }
      if (c.status === 'scheduled' && c.sendAt) {
        return `Scheduled for ${dayjs(c.sendAt).format('MMM D, YYYY HH:mm')}`;
      }
      return `${this.$t(`campaigns.status.${c.status}`)} · ${d.format('MMM D, YYYY HH:mm')}`;
    },

    audienceLists() {
      const attached = this.campaign.lists || [];
      const all = (this.lists && this.lists.results) || [];
      return attached.map((l) => {
        const full = all.find((x) => x.id === l.id);
        return {
          id: l.id,
          name: l.name || (full && full.name) || `List ${l.id}`,
          subscriberCount: (full && full.subscriberCount) || l.subscriberCount || 0,
        };
      });
    },

    audienceTotal() {
      return this.audienceLists.reduce((sum, l) => sum + (Number(l.subscriberCount) || 0), 0);
    },

    subscribersRoute() {
      const first = this.audienceLists[0];
      if (first) {
        return { name: 'subscribers_list', params: { listID: first.id } };
      }
      return { name: 'subscribers' };
    },
  },

  watch: {
    'campaign.id': {
      immediate: true,
      handler(id) {
        if (id) this.loadSummary(id);
      },
    },
  },

  methods: {
    n(num) {
      return this.$utils.formatNumber(Number(num) || 0);
    },

    pct(num, den) {
      if (!den || den <= 0) return '0%';
      const v = (Number(num) / den) * 100;
      if (!v) return '0%';
      return `${v.toFixed(2)}%`;
    },

    listShare(list) {
      return this.pct(list.subscriberCount || 0, this.audienceTotal);
    },

    listStat(list, campaignValue) {
      if (this.audienceLists.length <= 1) return Number(campaignValue) || 0;
      const total = this.audienceTotal;
      if (!total) return 0;
      return Math.round(((Number(list.subscriberCount) || 0) / total) * (Number(campaignValue) || 0));
    },

    loadSummary(id) {
      if (!this.$api.getCampaignReportSummary) return;
      this.$api.getCampaignReportSummary(id).then((data) => {
        const s = (data && data.summary) || {};
        this.summary = {
          uniqueOpens: Number(s.uniqueOpens) || 0,
          totalOpens: Number(s.totalOpens) || 0,
          uniqueClicks: Number(s.uniqueClicks) || 0,
          totalClicks: Number(s.totalClicks) || 0,
          unsubscribes: Number(s.unsubscribes) || 0,
          softBounces: Number(s.softBounces) || 0,
          hardBounces: Number(s.hardBounces) || 0,
          complaints: Number(s.complaints) || 0,
          bounces: Number(s.bounces) || 0,
        };
        this.bounceReasons = data.bounceReasons || [];
      }).catch(() => {});
    },

    onShare() {
      const url = window.location.href;
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(url).then(() => {
          this.$utils.toast('Report link copied');
        }).catch(() => {
          this.$utils.toast('Unable to copy link', 'is-danger');
        });
        return;
      }
      this.$utils.toast(url);
    },
  },

  mounted() {
    const { tab } = this.$route.query;
    if (tab && this.tabs.some((t) => t.id === tab)) {
      this.activeTab = tab;
    }
    if (!this.lists || !this.lists.results || this.lists.results.length === 0) {
      this.$api.getLists({ per_page: 'all', minimal: true }).catch(() => {});
    }
  },
});
</script>

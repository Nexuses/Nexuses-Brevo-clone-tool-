<template>
  <!-- eslint-disable vue/no-mutating-props -->
  <div class="campaign-brevo" :class="{ 'is-design': isDesignOpen }">
    <div v-if="isDesignOpen" class="campaign-brevo-design">
      <header class="campaign-brevo-design__bar">
        <div class="campaign-brevo-design__title">{{ form.name || 'Untitled campaign' }}</div>
        <div class="campaign-brevo-design__actions">
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--outline" @click="openPreviewTest">
            Preview &amp; Test
          </button>
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--dark" @click="saveAndQuitDesign">
            Save &amp; Quit
          </button>
        </div>
      </header>
      <div class="campaign-brevo-design__body">
        <editor
          v-if="campaignId"
          v-model="form.content"
          :id="campaignId"
          :title="form.name"
          :disabled="!canEdit"
          :templates="templates"
          :content-types="contentTypes"
        />
        <p v-else class="campaign-brevo__hint">Save recipients and subject first to open the editor.</p>
      </div>
    </div>

    <template v-else>
      <nav class="campaign-brevo__crumbs" aria-label="Breadcrumb">
        <router-link :to="{ name: 'campaigns' }">Email campaigns</router-link>
        <span>/</span>
        <span>Create an email campaign</span>
      </nav>

      <header class="campaign-brevo__header">
        <div class="campaign-brevo__title-wrap">
          <div v-if="!isEditingName" class="campaign-brevo__title-row">
            <router-link :to="{ name: 'campaigns' }" class="campaign-brevo__back" aria-label="Back to campaigns">
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
                <path d="M11.5 4.5L7 9l4.5 4.5" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"
                  stroke-linejoin="round" />
              </svg>
            </router-link>
            <h1 class="campaign-brevo__title">{{ displayName }}</h1>
            <button
              v-if="canEdit"
              type="button"
              class="campaign-brevo__icon-btn"
              aria-label="Rename campaign"
              @click="startRename"
            >
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none" aria-hidden="true">
                <path d="M8.7 2.6l3.7 3.7L4.8 14H1.2v-3.6L8.7 2.6z" stroke="currentColor" stroke-width="1.35"
                  stroke-linejoin="round" />
              </svg>
            </button>
            <span class="campaign-brevo__status">{{ statusLabel }}</span>
          </div>
          <div v-else class="campaign-brevo__title-row">
            <input
              ref="nameInput"
              v-model="nameDraft"
              class="campaign-brevo__name-input"
              maxlength="200"
              aria-label="Campaign name"
              @keydown.enter.prevent="saveName"
              @keydown.esc.prevent="isEditingName = false"
              @blur="saveName"
            />
          </div>
        </div>
        <div v-if="canManage || canSend" class="campaign-brevo__header-actions">
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--outline" @click="openPreviewTest">
            Preview &amp; Test
          </button>
          <button
            v-if="canUnSchedule"
            type="button"
            class="campaign-brevo__btn campaign-brevo__btn--outline"
            @click="$emit('unschedule')"
          >
            Unschedule
          </button>
          <button
            v-else-if="canSend"
            type="button"
            class="campaign-brevo__btn campaign-brevo__btn--dark"
            @click="openPanel('schedule')"
          >
            Schedule
          </button>
        </div>
      </header>

      <b-loading :active="loading" :is-full-page="false" />

      <div class="campaign-brevo__stack">
        <section class="campaign-brevo__card">
          <div class="campaign-brevo__card-top">
            <button type="button" class="campaign-brevo__languages" @click="onAddLanguages">
              Add languages
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
                <path d="M7 1.2l1.5 3.1 3.4.5-2.45 2.4.6 3.4L7 9l-3.05 1.6.6-3.4L2.1 4.8l3.4-.5L7 1.2z"
                  fill="#E8B931" stroke="#C99514" stroke-width="0.6" />
              </svg>
            </button>
          </div>

          <div class="campaign-brevo__row">
            <div class="campaign-brevo__row-main">
              <span class="campaign-brevo__check" :class="{ 'is-done': senderDone }" aria-hidden="true">
                <svg v-if="senderDone" width="11" height="11" viewBox="0 0 12 12" fill="none">
                  <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                    stroke-linejoin="round" />
                </svg>
              </span>
              <div>
                <h2>Sender</h2>
                <p>{{ senderSummary }}</p>
              </div>
            </div>
            <button type="button" class="campaign-brevo__btn campaign-brevo__btn--ghost" @click="openPanel('sender')">
              Manage sender
            </button>
          </div>

          <div class="campaign-brevo__row">
            <div class="campaign-brevo__row-main">
              <span class="campaign-brevo__check" :class="{ 'is-done': recipientsDone }" aria-hidden="true">
                <svg v-if="recipientsDone" width="11" height="11" viewBox="0 0 12 12" fill="none">
                  <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                    stroke-linejoin="round" />
                </svg>
              </span>
              <div>
                <h2>Recipients</h2>
                <p>{{ recipientsSummary }}</p>
              </div>
            </div>
            <button type="button" class="campaign-brevo__btn campaign-brevo__btn--ghost" @click="openPanel('recipients')">
              Edit recipients
            </button>
          </div>

          <div class="campaign-brevo__row">
            <div class="campaign-brevo__row-main">
              <span class="campaign-brevo__check" :class="{ 'is-done': subjectDone }" aria-hidden="true">
                <svg v-if="subjectDone" width="11" height="11" viewBox="0 0 12 12" fill="none">
                  <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                    stroke-linejoin="round" />
                </svg>
              </span>
              <div>
                <h2>Subject</h2>
                <p>{{ subjectSummary }}</p>
              </div>
            </div>
            <button type="button" class="campaign-brevo__btn campaign-brevo__btn--ghost" @click="openPanel('subject')">
              Edit subject
            </button>
          </div>

          <div class="campaign-brevo__row campaign-brevo__row--design">
            <div class="campaign-brevo__row-main">
              <span class="campaign-brevo__check" :class="{ 'is-done': designDone }" aria-hidden="true">
                <svg v-if="designDone" width="11" height="11" viewBox="0 0 12 12" fill="none">
                  <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                    stroke-linejoin="round" />
                </svg>
              </span>
              <div class="campaign-brevo__design-copy">
                <div class="campaign-brevo__design-title">
                  <h2>Design</h2>
                  <div class="campaign-brevo__dd" ref="designDd">
                    <button
                      type="button"
                      class="campaign-brevo__more"
                      aria-label="Design options"
                      @click.stop="designMenuOpen = !designMenuOpen"
                    >
                      <span /><span /><span />
                    </button>
                    <ul v-if="designMenuOpen" class="campaign-brevo__menu">
                      <li><button type="button" @click="resetDesign">Reset design</button></li>
                      <li><button type="button" @click="viewHtml">View the HTML code</button></li>
                      <li><button type="button" @click="downloadHtml">Download the HTML code</button></li>
                      <li><button type="button" @click="viewPlain">View the plain text version</button></li>
                    </ul>
                  </div>
                </div>
                <div class="campaign-brevo__thumb">
                  <iframe
                    v-if="campaignId && designDone"
                    class="campaign-brevo__thumb-frame"
                    :src="previewSrc"
                    title="Email preview"
                    tabindex="-1"
                    sandbox="allow-scripts"
                  />
                  <div v-else class="campaign-brevo__thumb-empty">No design yet</div>
                </div>
              </div>
            </div>
            <button type="button" class="campaign-brevo__btn campaign-brevo__btn--ghost" @click="openDesign">
              Edit design
            </button>
          </div>
        </section>

        <section class="campaign-brevo__card">
          <div class="campaign-brevo__row">
            <div class="campaign-brevo__row-main">
              <span class="campaign-brevo__check is-muted" aria-hidden="true" />
              <div>
                <h2>Additional settings</h2>
                <p v-if="settingsSummary">{{ settingsSummary }}</p>
                <p v-else class="is-muted">Sending and tracking</p>
              </div>
            </div>
            <button type="button" class="campaign-brevo__btn campaign-brevo__btn--ghost" @click="openPanel('settings')">
              Edit settings
            </button>
          </div>
        </section>
      </div>
    </template>

    <!-- Sender modal -->
    <div v-if="panel === 'sender'" class="campaign-brevo-modal" role="dialog" aria-modal="true" aria-labelledby="sender-title">
      <button type="button" class="campaign-brevo-modal__backdrop" aria-label="Close" @click="closePanel" />
      <div class="campaign-brevo-modal__card">
        <header class="campaign-brevo-modal__head">
          <div class="campaign-brevo-modal__head-main">
            <span class="campaign-brevo__check is-done" aria-hidden="true">
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                  stroke-linejoin="round" />
              </svg>
            </span>
            <div>
              <h2 id="sender-title">Sender</h2>
              <p>Who is sending this email campaign?</p>
            </div>
          </div>
          <button type="button" class="campaign-brevo__icon-btn" aria-label="Close" @click="closePanel">×</button>
        </header>
        <div class="campaign-brevo-modal__body campaign-brevo-modal__split">
          <div>
            <label class="campaign-brevo__field">
              <span>Email address</span>
              <select v-model="draft.senderEmail" :disabled="!canEdit">
                <option v-for="em in senderEmails" :key="em" :value="em">{{ em }}</option>
              </select>
            </label>
            <label class="campaign-brevo__field">
              <span>Name</span>
              <input v-model="draft.senderName" maxlength="200" :disabled="!canEdit" placeholder="Sender name" />
            </label>
          </div>
          <div class="campaign-brevo__phone" aria-hidden="true">
            <div class="campaign-brevo__phone-screen">
              <div class="campaign-brevo__inbox-row is-active">
                <strong>{{ draft.senderName || 'Sender name' }}</strong>
                <span>17:45</span>
                <em>{{ form.subject || 'Message subject...' }}</em>
                <small>{{ form.previewText || 'Your preview text' }}</small>
              </div>
              <div class="campaign-brevo__inbox-row" />
              <div class="campaign-brevo__inbox-row" />
            </div>
            <p>Actual email preview may vary depending on the email client.</p>
          </div>
        </div>
        <footer class="campaign-brevo-modal__foot">
          <button type="button" class="campaign-brevo__link" @click="closePanel">Cancel</button>
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--dark" @click="saveSender">Save</button>
        </footer>
      </div>
    </div>

    <!-- Recipients modal -->
    <div v-if="panel === 'recipients'" class="campaign-brevo-modal" role="dialog" aria-modal="true" aria-labelledby="recip-title">
      <button type="button" class="campaign-brevo-modal__backdrop" aria-label="Close" @click="closePanel" />
      <div class="campaign-brevo-modal__card is-wide">
        <header class="campaign-brevo-modal__head">
          <div class="campaign-brevo-modal__head-main">
            <span class="campaign-brevo__check" :class="{ 'is-done': recipientsDone }" aria-hidden="true">
              <svg v-if="recipientsDone" width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                  stroke-linejoin="round" />
              </svg>
            </span>
            <div>
              <h2 id="recip-title">Recipients</h2>
              <p>{{ recipientsSummary }}</p>
            </div>
          </div>
          <button type="button" class="campaign-brevo__icon-btn" aria-label="Close" @click="closePanel">×</button>
        </header>
        <div class="campaign-brevo-modal__body">
          <label class="campaign-brevo__field">
            <span>Send to</span>
            <div class="campaign-brevo__chips">
              <span v-for="l in draft.lists" :key="l.id" class="campaign-brevo__chip">
                {{ l.name }}
                <button v-if="canEdit" type="button" aria-label="Remove list" @click="removeDraftList(l.id)">×</button>
              </span>
              <select v-if="canEdit" :value="''" aria-label="Add list" @change="addDraftList($event)">
                <option value="">Select list(s), segment(s) or individual contacts</option>
                <option v-for="l in availableLists" :key="l.id" :value="l.id">{{ l.name }}</option>
              </select>
            </div>
          </label>
          <label class="campaign-brevo__checkline">
            <input v-model="draft.skipUnengaged" type="checkbox" :disabled="!canEdit" />
            Don’t send to unengaged contacts
          </label>
          <button type="button" class="campaign-brevo__advanced" @click="advancedOpen = !advancedOpen">
            Advanced options
            <span>{{ advancedOpen ? '▴' : '▾' }}</span>
          </button>
          <div v-if="advancedOpen">
            <label class="campaign-brevo__field">
              <span>Don’t send to</span>
              <select disabled>
                <option>Select list(s), segment(s) or individual contacts</option>
              </select>
            </label>
            <p class="campaign-brevo__hint">Exclude lists or segments after you save this campaign.</p>
          </div>
          <div class="campaign-brevo__recip-foot">
            <strong>{{ recipientCountLabel }}</strong>
            <p>Send to as many recipients as you wish, within your plan limits.</p>
          </div>
        </div>
        <footer class="campaign-brevo-modal__foot">
          <button type="button" class="campaign-brevo__link" @click="closePanel">Cancel</button>
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--dark" @click="saveRecipients">Save</button>
        </footer>
      </div>
    </div>

    <!-- Subject modal -->
    <div v-if="panel === 'subject'" class="campaign-brevo-modal" role="dialog" aria-modal="true" aria-labelledby="subj-title">
      <button type="button" class="campaign-brevo-modal__backdrop" aria-label="Close" @click="closePanel" />
      <div class="campaign-brevo-modal__card is-wide">
        <header class="campaign-brevo-modal__head">
          <div class="campaign-brevo-modal__head-main">
            <span class="campaign-brevo__check" :class="{ 'is-done': subjectDone }" aria-hidden="true">
              <svg v-if="subjectDone" width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path d="M2.5 6.2l2.4 2.4 4.6-5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"
                  stroke-linejoin="round" />
              </svg>
            </span>
            <div>
              <h2 id="subj-title">Subject</h2>
              <p>Add a subject line for this campaign.</p>
            </div>
          </div>
          <button type="button" class="campaign-brevo__icon-btn" aria-label="Close" @click="closePanel">×</button>
        </header>
        <div class="campaign-brevo-modal__body campaign-brevo-modal__split">
          <div>
            <label class="campaign-brevo__field">
              <span>Subject line *</span>
              <div class="campaign-brevo__textarea-wrap">
                <textarea
                  ref="subjectBox"
                  v-model="draft.subject"
                  rows="3"
                  maxlength="5000"
                  :disabled="!canEdit"
                  required
                />
                <div class="campaign-brevo__toolbox">
                  <button type="button" aria-label="Insert emoji" @click="insertAt('subject', '🎉')">☺</button>
                  <button type="button" aria-label="Insert personalization" @click="insertName('subject')">{}</button>
                  <button type="button" aria-label="Suggest subject" @click="suggestSubject">✦</button>
                </div>
              </div>
            </label>
            <label class="campaign-brevo__field">
              <span>Preview text</span>
              <div class="campaign-brevo__textarea-wrap">
                <textarea v-model="draft.previewText" rows="3" maxlength="500" :disabled="!canEdit" />
                <div class="campaign-brevo__toolbox">
                  <button type="button" aria-label="Insert emoji" @click="insertAt('previewText', '✨')">☺</button>
                  <button type="button" aria-label="Insert personalization" @click="insertName('previewText')">{}</button>
                </div>
              </div>
            </label>
          </div>
          <div class="campaign-brevo__phone" aria-hidden="true">
            <div class="campaign-brevo__phone-screen">
              <div class="campaign-brevo__inbox-row is-active">
                <strong>{{ senderName || 'Sender' }}</strong>
                <span>17:45</span>
                <em>{{ draft.subject || 'Message subject...' }}</em>
                <small>{{ draft.previewText || 'Your preview text' }}</small>
              </div>
              <div class="campaign-brevo__inbox-row" />
              <div class="campaign-brevo__inbox-row" />
            </div>
            <p>Actual email preview may vary depending on the email client.</p>
          </div>
        </div>
        <footer class="campaign-brevo-modal__foot">
          <button type="button" class="campaign-brevo__link" @click="closePanel">Cancel</button>
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--dark" @click="saveSubject">Save</button>
        </footer>
      </div>
    </div>

    <!-- Additional settings modal -->
    <div v-if="panel === 'settings'" class="campaign-brevo-modal" role="dialog" aria-modal="true" aria-labelledby="set-title">
      <button type="button" class="campaign-brevo-modal__backdrop" aria-label="Close" @click="closePanel" />
      <div class="campaign-brevo-modal__card">
        <header class="campaign-brevo-modal__head">
          <h2 id="set-title">Additional settings</h2>
          <button type="button" class="campaign-brevo__icon-btn" aria-label="Close" @click="closePanel">×</button>
        </header>
        <div class="campaign-brevo-modal__body">
          <h3 class="campaign-brevo__section-label">Personalization</h3>
          <label class="campaign-brevo__toggle">
            <span>Personalize the ‘Send To’ field</span>
            <input v-model="draft.personalizeTo" type="checkbox" :disabled="!canEdit" />
          </label>

          <h3 class="campaign-brevo__section-label">Sending and Tracking</h3>
          <label class="campaign-brevo__toggle">
            <span>Use a different Reply-to address</span>
            <input v-model="draft.useReplyTo" type="checkbox" :disabled="!canEdit" />
          </label>
          <label v-if="draft.useReplyTo" class="campaign-brevo__field">
            <span class="is-sr-only">Reply-to email</span>
            <input v-model="draft.replyTo" type="email" :disabled="!canEdit" placeholder="reply@example.com" />
          </label>
          <label class="campaign-brevo__toggle">
            <span>Activate UTM tracking</span>
            <input v-model="draft.utmOn" type="checkbox" :disabled="!canEdit" />
          </label>
          <div v-if="draft.utmOn" class="campaign-brevo__utm">
            <input v-model="draft.utmSource" :disabled="!canEdit" placeholder="utm_source" />
            <input v-model="draft.utmMedium" :disabled="!canEdit" placeholder="utm_medium" />
            <input v-model="draft.utmCampaign" :disabled="!canEdit" placeholder="utm_campaign" />
          </div>
          <label class="campaign-brevo__toggle">
            <span>Add an attachment</span>
            <input v-model="draft.attachOn" type="checkbox" :disabled="!canEdit" />
          </label>
          <div v-if="draft.attachOn" class="campaign-brevo__attach">
            <button type="button" class="campaign-brevo__btn campaign-brevo__btn--ghost" @click="$emit('attach')">
              Choose files
            </button>
            <span v-if="form.media && form.media.length">{{ form.media.length }} attached</span>
          </div>
          <label class="campaign-brevo__toggle">
            <span>Add a tag</span>
            <input v-model="draft.tagsOn" type="checkbox" :disabled="!canEdit" />
          </label>
          <div v-if="draft.tagsOn">
            <b-taginput v-model="draft.tags" :disabled="!canEdit" ellipsis icon="tag-outline" placeholder="Tags" />
          </div>
        </div>
        <footer class="campaign-brevo-modal__foot">
          <button type="button" class="campaign-brevo__link" @click="closePanel">Cancel</button>
          <button type="button" class="campaign-brevo__btn campaign-brevo__btn--dark" @click="saveSettings">Save</button>
        </footer>
      </div>
    </div>

    <!-- Schedule modal -->
    <div v-if="panel === 'schedule'" class="campaign-brevo-modal" role="dialog" aria-modal="true" aria-labelledby="sched-title">
      <button type="button" class="campaign-brevo-modal__backdrop" aria-label="Close" @click="closePanel" />
      <div class="campaign-brevo-modal__card">
        <header class="campaign-brevo-modal__head">
          <h2 id="sched-title">Schedule</h2>
          <button type="button" class="campaign-brevo__icon-btn" aria-label="Close" @click="closePanel">×</button>
        </header>
        <div class="campaign-brevo-modal__body">
          <p v-if="!canScheduleSend" class="campaign-brevo__hint">
            Complete sender, recipients, and subject before scheduling.
          </p>
          <label class="campaign-brevo__radio">
            <input v-model="draft.sendNow" type="radio" :value="true" />
            Send now
          </label>
          <label class="campaign-brevo__radio">
            <input v-model="draft.sendNow" type="radio" :value="false" />
            Schedule for later
          </label>
          <b-datetimepicker
            v-if="!draft.sendNow"
            v-model="draft.sendAt"
            required
            editable
            mobile-native
            icon="calendar-clock"
            :timepicker="{ hourFormat: '24' }"
            placeholder="Date and time"
          />
        </div>
        <footer class="campaign-brevo-modal__foot">
          <button type="button" class="campaign-brevo__link" @click="closePanel">Cancel</button>
          <button
            type="button"
            class="campaign-brevo__btn campaign-brevo__btn--dark"
            :disabled="!canScheduleSend"
            @click="confirmSchedule"
          >
            {{ draft.sendNow ? 'Send' : 'Schedule' }}
          </button>
        </footer>
      </div>
    </div>

    <!-- Preview & Test modal -->
    <div v-if="panel === 'previewTest'" class="campaign-brevo-modal" role="dialog" aria-modal="true" aria-labelledby="pt-title">
      <button type="button" class="campaign-brevo-modal__backdrop" aria-label="Close" @click="closePanel" />
      <div class="campaign-brevo-modal__card is-wide is-preview-test">
        <header class="campaign-brevo-modal__head">
          <h2 id="pt-title">Preview &amp; Test</h2>
          <button type="button" class="campaign-brevo__icon-btn" aria-label="Close" @click="closePanel">×</button>
        </header>
        <div class="campaign-brevo-modal__body campaign-brevo-modal__split is-preview-test">
          <div class="campaign-brevo-modal__test">
            <p class="campaign-brevo__hint">Send a test to check rendering before scheduling.</p>
            <b-taginput
              v-model="form.testEmails"
              :before-adding="$utils.validateEmail"
              :disabled="isNew"
              ellipsis
              icon="email-outline"
              placeholder="Test email addresses"
            />
            <button
              type="button"
              class="campaign-brevo__btn campaign-brevo__btn--dark mt-3"
              :disabled="isNew || !(form.testEmails && form.testEmails.length)"
              @click="$emit('test')"
            >
              Send test
            </button>
            <p v-if="isNew" class="campaign-brevo__hint mt-2">Save the campaign first to send a test.</p>
          </div>
          <div v-if="campaignId" class="campaign-brevo-modal__preview">
            <iframe
              class="campaign-brevo-modal__preview-frame"
              :src="previewSrc"
              title="Campaign preview"
              sandbox="allow-scripts"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
/* eslint-disable vue/no-mutating-props */
import Vue from 'vue';
import Editor from './Editor.vue';
import { uris } from '../constants';

export default Vue.extend({
  name: 'CampaignSetup',

  components: { Editor },

  props: {
    form: { type: Object, required: true },
    data: { type: Object, default: () => ({}) },
    isNew: { type: Boolean, default: false },
    canEdit: { type: Boolean, default: true },
    canManage: { type: Boolean, default: false },
    canSend: { type: Boolean, default: false },
    canUnSchedule: { type: Boolean, default: false },
    loading: { type: Boolean, default: false },
    lists: { type: Array, default: () => [] },
    templates: { type: Array, default: () => [] },
    contentTypes: { type: Object, default: () => ({}) },
    fromEmailDefault: { type: String, default: '' },
  },

  data() {
    return {
      panel: null,
      isDesignOpen: false,
      designMenuOpen: false,
      isEditingName: false,
      nameDraft: '',
      advancedOpen: false,
      draft: this.emptyDraft(),
    };
  },

  computed: {
    campaignId() {
      return this.data && this.data.id ? this.data.id : 0;
    },

    displayName() {
      return this.form.name || 'Untitled campaign';
    },

    statusLabel() {
      const s = (this.data && this.data.status) || 'draft';
      if (s === 'draft') return 'Draft';
      if (s === 'scheduled') return 'Scheduled';
      if (s === 'paused') return 'Paused';
      return s;
    },

    parsedFrom() {
      return this.parseFrom(this.form.fromEmail || this.fromEmailDefault || '');
    },

    senderName() {
      return this.parsedFrom.name || '';
    },

    senderEmail() {
      return this.parsedFrom.email || '';
    },

    senderEmails() {
      const out = [];
      const def = this.parseFrom(this.fromEmailDefault || '').email;
      if (def) out.push(def);
      if (this.senderEmail && out.indexOf(this.senderEmail) < 0) out.push(this.senderEmail);
      if (!out.length) out.push('noreply@example.com');
      return out;
    },

    senderDone() {
      return !!(this.form.fromEmail || this.fromEmailDefault);
    },

    senderSummary() {
      const { name, email } = this.parsedFrom;
      if (name && email) return `${name} • ${email}`;
      return email || name || 'Add a sender';
    },

    recipientsDone() {
      return !!(this.form.lists && this.form.lists.length);
    },

    recipientCount() {
      return (this.form.lists || []).reduce((sum, l) => sum + (Number(l.subscriberCount) || 0), 0);
    },

    recipientCountLabel() {
      const n = this.$utils.formatNumber(this.recipientCount);
      const lists = (this.form.lists || []).length;
      if (!lists) return '0 recipients';
      if (this.recipientCount > 0) return `${n} recipients`;
      return `${lists} list${lists === 1 ? '' : 's'} selected`;
    },

    recipientsSummary() {
      if (!this.recipientsDone) return 'Select who should receive this campaign';
      return this.recipientCountLabel;
    },

    subjectDone() {
      return !!(this.form.subject && this.form.subject.trim());
    },

    subjectSummary() {
      if (!this.subjectDone) return 'Add a subject line for this campaign';
      return `Subject: ${this.form.subject}`;
    },

    designDone() {
      const body = this.form.content && this.form.content.body;
      return !!(body && String(body).trim());
    },

    previewSrc() {
      if (!this.campaignId) return 'about:blank';
      return uris.previewCampaign.replace(':id', this.campaignId);
    },

    settingsSummary() {
      const bits = [];
      if (this.form.replyTo) bits.push("You're using a different 'Reply-to' Email address.");
      if (this.form.tracking && (this.form.tracking.utm_source || this.form.tracking.utm_campaign)) {
        bits.push('UTM tracking is on.');
      }
      if (this.form.media && this.form.media.length) bits.push(`${this.form.media.length} attachment(s).`);
      return bits.join(' ');
    },

    availableLists() {
      const selected = {};
      (this.draft.lists || []).forEach((l) => { selected[l.id] = true; });
      return (this.lists || []).filter((l) => !selected[l.id]);
    },

    canScheduleSend() {
      return this.senderDone && this.recipientsDone && this.subjectDone;
    },
  },

  mounted() {
    document.addEventListener('click', this.onDocClick);
  },

  beforeDestroy() {
    document.removeEventListener('click', this.onDocClick);
  },

  methods: {
    emptyDraft() {
      return {
        senderName: '',
        senderEmail: '',
        lists: [],
        skipUnengaged: false,
        subject: '',
        previewText: '',
        personalizeTo: false,
        useReplyTo: false,
        replyTo: '',
        utmOn: false,
        utmSource: '',
        utmMedium: '',
        utmCampaign: '',
        attachOn: false,
        tagsOn: false,
        tags: [],
        sendNow: true,
        sendAt: null,
      };
    },

    parseFrom(raw) {
      const s = String(raw || '').trim();
      const m = s.match(/^(.*?)\s*<([^>]+)>$/);
      if (m) {
        return { name: m[1].replace(/^"|"$/g, '').trim(), email: m[2].trim() };
      }
      if (s.indexOf('@') > -1) {
        const local = s.split('@')[0].replace(/[._]/g, ' ');
        const name = local.replace(/\b\w/g, (c) => c.toUpperCase());
        return { name, email: s };
      }
      return { name: s, email: '' };
    },

    formatFrom(name, email) {
      const em = (email || '').trim();
      const nm = (name || '').trim();
      if (nm && em) return `${nm} <${em}>`;
      return em || nm;
    },

    onDocClick() {
      this.designMenuOpen = false;
    },

    openPanel(name) {
      this.designMenuOpen = false;
      this.draft = {
        ...this.emptyDraft(),
        senderName: this.senderName,
        senderEmail: this.senderEmail || this.senderEmails[0],
        lists: [...(this.form.lists || [])],
        subject: this.form.subject || '',
        previewText: this.form.previewText || '',
        useReplyTo: !!(this.form.replyTo && this.form.replyTo.trim()),
        replyTo: this.form.replyTo || '',
        utmOn: !!(this.form.tracking && (this.form.tracking.utm_source || this.form.tracking.utm_campaign)),
        utmSource: (this.form.tracking && this.form.tracking.utm_source) || '',
        utmMedium: (this.form.tracking && this.form.tracking.utm_medium) || '',
        utmCampaign: (this.form.tracking && this.form.tracking.utm_campaign) || '',
        attachOn: !!(this.form.media && this.form.media.length),
        tagsOn: !!(this.form.tags && this.form.tags.length),
        tags: [...(this.form.tags || [])],
        sendNow: true,
        sendAt: this.form.sendAtDate || null,
      };
      this.panel = name;
      this.advancedOpen = false;
    },

    closePanel() {
      this.panel = null;
    },

    persist(thenCreate) {
      this.$emit('save', { createIfNew: !!thenCreate });
    },

    startRename() {
      this.nameDraft = this.form.name || '';
      this.isEditingName = true;
      this.$nextTick(() => {
        if (this.$refs.nameInput) this.$refs.nameInput.focus();
      });
    },

    saveName() {
      if (!this.isEditingName) return;
      this.isEditingName = false;
      const name = (this.nameDraft || '').trim() || 'Untitled campaign';
      this.form.name = name;
      if (!this.isNew) this.persist(false);
    },

    saveSender() {
      this.form.fromEmail = this.formatFrom(this.draft.senderName, this.draft.senderEmail);
      this.closePanel();
      this.persist(false);
    },

    addDraftList(ev) {
      const el = ev && ev.target;
      const id = parseInt(el && el.value, 10);
      if (el) el.selectedIndex = 0;
      if (!id) return;
      const found = (this.lists || []).find((l) => l.id === id);
      if (found) this.draft.lists.push(found);
    },

    removeDraftList(id) {
      this.draft.lists = this.draft.lists.filter((l) => l.id !== id);
    },

    saveRecipients() {
      this.form.lists = [...this.draft.lists];
      this.closePanel();
      this.persist(true);
    },

    insertAt(field, token) {
      const cur = this.draft[field] || '';
      this.draft[field] = `${cur}${token}`;
    },

    insertName(field) {
      this.insertAt(field, '{{ .Subscriber.Name }}');
    },

    suggestSubject() {
      if (this.draft.subject && this.draft.subject.trim()) return;
      this.draft.subject = this.form.name || 'Your latest update';
    },

    saveSubject() {
      const sub = (this.draft.subject || '').trim();
      if (!sub) {
        this.$utils.toast('Add a subject line.');
        return;
      }
      this.form.subject = sub;
      this.form.previewText = this.draft.previewText || '';
      if (!this.form.name) this.form.name = sub.slice(0, 80);
      this.closePanel();
      this.persist(true);
    },

    saveSettings() {
      this.form.replyTo = this.draft.useReplyTo ? (this.draft.replyTo || '') : '';
      if (!this.form.tracking) this.form.tracking = {};
      if (this.draft.utmOn) {
        this.form.tracking.utm_source = this.draft.utmSource || 'listmonk';
        this.form.tracking.utm_medium = this.draft.utmMedium || 'email';
        this.form.tracking.utm_campaign = this.draft.utmCampaign || (this.form.name || 'campaign');
        this.form.tracking.enabled = true;
      } else {
        this.form.tracking.utm_source = '';
        this.form.tracking.utm_medium = '';
        this.form.tracking.utm_campaign = '';
      }
      if (this.draft.tagsOn) this.form.tags = [...(this.draft.tags || [])];
      this.closePanel();
      this.persist(false);
    },

    openDesign() {
      this.designMenuOpen = false;
      if (this.isNew) {
        this.$emit('ensure', () => {
          this.isDesignOpen = true;
          this.$emit('design-open', true);
        });
        return;
      }
      this.isDesignOpen = true;
      this.$emit('design-open', true);
    },

    saveAndQuitDesign() {
      this.isDesignOpen = false;
      this.$emit('design-open', false);
      this.persist(false);
    },

    resetDesign() {
      this.designMenuOpen = false;
      this.form.content.body = '';
      this.form.content.bodySource = null;
      this.persist(false);
    },

    viewHtml() {
      this.designMenuOpen = false;
      this.form.content.contentType = 'html';
      this.openDesign();
    },

    downloadHtml() {
      this.designMenuOpen = false;
      const html = this.form.content.body || '';
      const blob = new Blob([html], { type: 'text/html' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${(this.form.name || 'campaign').replace(/\s+/g, '-')}.html`;
      a.click();
      URL.revokeObjectURL(url);
    },

    viewPlain() {
      this.designMenuOpen = false;
      this.form.content.contentType = 'plain';
      this.openDesign();
    },

    onAddLanguages() {
      this.$utils.toast('Multiple languages are not available in this plan');
    },

    openPreviewTest() {
      this.panel = 'previewTest';
    },

    confirmSchedule() {
      if (!this.canScheduleSend) return;
      if (!this.draft.sendNow && !this.draft.sendAt) {
        this.$utils.toast('Pick a date and time.');
        return;
      }
      this.closePanel();
      this.$emit('schedule', {
        sendNow: this.draft.sendNow,
        at: this.draft.sendAt,
      });
    },
  },
});
</script>

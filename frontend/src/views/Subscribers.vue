<template>
  <section class="subscribers contacts-brevo">
    <header class="contacts-brevo__header">
      <h1 class="contacts-brevo__title">
        <template v-if="currentList">{{ currentList.name }}</template>
        <template v-else>Contacts</template>
      </h1>
      <div class="contacts-brevo__actions">
        <button
          v-if="$can('subscribers:manage')"
          type="button"
          class="contacts-brevo__btn contacts-brevo__btn--outline"
          @click="showNewForm"
          data-cy="btn-new"
        >
          {{ currentList ? 'Add a contact' : 'Create a contact' }}
        </button>
        <router-link
          v-if="$can('subscribers:import')"
          :to="importRoute"
          class="contacts-brevo__btn contacts-brevo__btn--dark"
        >
          Import contacts
        </router-link>
      </div>
    </header>

    <div v-if="currentList" class="contacts-brevo__list-actions">
      <button
        v-if="$can('subscribers:manage')"
        type="button"
        class="contacts-brevo__chip"
        @click="showNewForm"
      >
        <b-icon icon="plus" size="is-small" />
        Add user to this list
      </button>
      <router-link
        v-if="$can('subscribers:import')"
        :to="importRoute"
        class="contacts-brevo__chip"
      >
        <b-icon icon="file-upload-outline" size="is-small" />
        Import users to this list
      </router-link>
    </div>

    <div class="contacts-brevo__views">
      <router-link
        :to="{ name: 'subscribers' }"
        class="contacts-brevo__view-tab"
        :class="{ 'is-active': !currentList }"
      >
        All contacts
      </router-link>
      <button
        v-if="currentList"
        type="button"
        class="contacts-brevo__view-tab is-active"
      >
        {{ currentList.name }}
      </button>
      <button
        v-if="$can('subscribers:manage')"
        type="button"
        class="contacts-brevo__view-add"
        aria-label="Create a contact"
        @click="showNewForm"
      >
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
          <path d="M7 2v10M2 7h10" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
        </svg>
      </button>
    </div>

    <div class="contacts-brevo__filter-bar">
      <div class="contacts-brevo__filter-left">
        <b-dropdown @change="onFilterListSelect">
          <template #trigger>
            <button type="button" class="contacts-brevo__chip">
              Load a list or a segment
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </button>
          </template>
          <b-dropdown-item v-if="currentList" :value="null">
            All contacts
          </b-dropdown-item>
          <b-dropdown-item
            v-for="l in availableLists"
            :key="l.id"
            :value="l.id"
          >
            {{ l.name }}
          </b-dropdown-item>
          <b-dropdown-item v-if="!availableLists.length" disabled>
            No lists available
          </b-dropdown-item>
        </b-dropdown>

        <button
          type="button"
          class="contacts-brevo__chip"
          :class="{ 'is-active': isFilterOpen || filters.length > 0 }"
          @click.prevent="toggleFilterPanel"
        >
          Add filter
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
      </div>
      <button type="button" class="contacts-brevo__gear" aria-label="Settings">
        <b-icon icon="cog-outline" size="is-small" />
      </button>
    </div>

    <div class="contacts-brevo__meta">
      <span class="contacts-brevo__meta-count">
        {{ $utils.formatNumber(subscribers.total || 0) }} contacts
        <span class="contacts-brevo__info" title="Total contacts matching this view">i</span>
      </span>
      <div class="contacts-brevo__meta-right">
        <button type="button" class="contacts-brevo__customize" @click="openColumnsDrawer">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M4 5h4v14H4V5zM10 5h4v14h-4V5zM16 5h4v14h-4V5z" stroke="currentColor" stroke-width="1.8" />
          </svg>
          Customize columns
        </button>
        <form class="contacts-brevo__search" @submit.prevent="onSubmit">
          <b-field>
            <b-input
              @input="onSimpleQueryInput"
              v-model="queryInput"
              expanded
              placeholder="Search"
              icon="magnify"
              ref="query"
              data-cy="search"
            />
          </b-field>
        </form>
      </div>
    </div>

    <div v-if="isFilterOpen" class="contacts-filter">
      <div v-for="(f, idx) in filters" :key="f.id" class="contacts-filter__block">
        <p v-if="idx > 0" class="contacts-filter__join-label">{{ f.join === 'or' ? 'Or' : 'And' }}</p>
        <div class="contacts-filter__row">
          <select v-model="f.field" class="contacts-filter__select" aria-label="Filter attribute" @change="onFilterFieldChange(f)">
            <option v-for="field in filterFieldDefs" :key="field.key" :value="field.key">{{ field.label }}</option>
          </select>
          <select v-model="f.op" class="contacts-filter__select" aria-label="Filter operator">
            <option v-for="op in operatorsFor(f)" :key="op.value" :value="op.value">{{ op.label }}</option>
          </select>
          <select
            v-if="filterType(f) === 'list' && !isEmptyOp(f.op)"
            v-model="f.value"
            class="contacts-filter__select contacts-filter__select--wide"
            aria-label="List"
          >
            <option value="">Select a list</option>
            <option v-for="l in availableLists" :key="l.id" :value="String(l.id)">{{ l.name }}</option>
          </select>
          <select
            v-else-if="filterType(f) === 'status' && !isEmptyOp(f.op)"
            v-model="f.value"
            class="contacts-filter__select"
            aria-label="Subscription status"
          >
            <option value="subscribed">Subscribed</option>
            <option value="blocklisted">Blocklisted</option>
          </select>
          <input
            v-else-if="filterType(f) === 'date' && f.op === 'last_n_days'"
            v-model="f.value"
            class="contacts-filter__input contacts-filter__input--sm"
            type="number"
            min="1"
            placeholder="Days"
            aria-label="Number of days"
          />
          <input
            v-else-if="filterType(f) === 'date' && !isEmptyOp(f.op)"
            v-model="f.value"
            class="contacts-filter__input"
            type="date"
            aria-label="Date"
          />
          <input
            v-else-if="filterType(f) === 'text' && !isEmptyOp(f.op)"
            v-model="f.value"
            class="contacts-filter__input"
            type="text"
            :placeholder="'Enter ' + fieldLabel(f.field)"
            aria-label="Filter value"
            @keydown.enter.prevent="applyFilters"
          />
          <button type="button" class="contacts-filter__remove" aria-label="Remove filter" @click="removeFilter(idx)">
            ×
          </button>
        </div>
      </div>
      <div class="contacts-filter__toolbar">
        <button type="button" class="contacts-filter__link" @click="addFilter('and')">+ And</button>
        <button type="button" class="contacts-filter__link" @click="addFilter('or')">+ Or</button>
        <div class="contacts-filter__actions">
          <button type="button" class="contacts-filter__clear" @click="clearFilters">Clear</button>
          <button type="button" class="contacts-filter__search" @click="applyFilters">Search contacts</button>
        </div>
      </div>
    </div>

    <b-table :data="subscribers.results ?? []" :loading="loading.subscribers" @check-all="onTableCheck"
      @check="onTableCheck" :checked-rows.sync="bulk.checked" :paginated="false"
      hoverable checkable backend-sorting @sort="onSort"
      class="contacts-brevo__table">
      <template #top-left>
        <div class="actions">
          <a class="a" href="#" @click.prevent="exportSubscribers" data-cy="btn-export-subscribers">
            <b-icon icon="cloud-download-outline" size="is-small" />
            {{ $t('subscribers.export') }}
          </a>
          <template v-if="bulk.checked.length > 0">
            <a class="a" href="#" @click.prevent="showBulkListForm" data-cy="btn-manage-lists">
              <b-icon icon="format-list-bulleted-square" size="is-small" /> Manage lists
            </a>
            <a class="a" href="#" @click.prevent="deleteSubscribers" data-cy="btn-delete-subscribers">
              <b-icon icon="trash-can-outline" size="is-small" /> Delete
            </a>
            <a class="a" href="#" @click.prevent="blocklistSubscribers" data-cy="btn-manage-blocklist">
              <b-icon icon="account-off-outline" size="is-small" /> Blocklist
            </a>
            <span class="a">
              {{ $t('globals.messages.numSelected', { num: numSelectedSubscribers }) }}
            </span>
          </template>
        </div>
      </template>

      <b-table-column v-slot="props" field="name" label="CONTACT" header-class="cy-email" sortable
        :td-attrs="$utils.tdID">
        <a :href="`/contacts/${props.row.id}`" @click.prevent="showEditForm(props.row)"
          class="contacts-brevo__name"
          :class="{ 'blocklisted': props.row.status === 'blocklisted' }">
          {{ contactDisplayName(props.row) }}
        </a>
      </b-table-column>

      <b-table-column
        v-for="col in visibleColumnDefs"
        :key="col.key"
        v-slot="props"
        :field="col.sortField || col.key"
        :label="col.label"
        :sortable="!!col.sortField"
      >
        <span
          v-if="col.key === 'subscribed'"
          class="contacts-brevo__sub-pill"
          :class="{ 'is-blocklisted': props.row.status === 'blocklisted' }"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M4 6h16v12H4V6z" stroke="currentColor" stroke-width="1.8" />
            <path d="M4 7l8 6 8-6" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" />
          </svg>
          Email
        </span>
        <span v-else-if="col.key === 'blocklisted'" class="contacts-brevo__cell">
          {{ props.row.status === 'blocklisted' ? 'Yes' : 'No' }}
        </span>
        <span v-else class="contacts-brevo__cell" :title="columnDisplay(props.row, col.key)">
          {{ columnDisplay(props.row, col.key) }}
        </span>
      </b-table-column>

      <b-table-column v-slot="props" cell-class="actions" align="right" width="90">
        <div class="contacts-brevo__row-actions">
          <a v-if="$can('subscribers:manage')" :href="`/contacts/${props.row.id}`"
            @click.prevent="showEditForm(props.row)" data-cy="btn-edit" :aria-label="$t('globals.buttons.edit')">
            <b-icon icon="pencil-outline" size="is-small" />
          </a>
          <a v-if="$can('subscribers:manage')" href="#" @click.prevent="deleteSubscriber(props.row)"
            data-cy="btn-delete" :aria-label="$t('globals.buttons.delete')">
            <b-icon icon="trash-can-outline" size="is-small" />
          </a>
        </div>
      </b-table-column>

      <template #empty v-if="!loading.subscribers">
        <empty-placeholder />
      </template>
    </b-table>

    <div class="contacts-brevo__pager">
      <div class="contacts-brevo__pager-meta">
        <span>{{ pagerSummary }}</span>
        <span class="contacts-brevo__pager-pages">Page {{ queryParams.page }} of {{ totalPages }}</span>
      </div>
      <div class="contacts-brevo__pager-controls">
        <label class="contacts-brevo__pager-size">
          Rows
          <select :value="queryParams.perPage" aria-label="Rows per page" @change="onPerPageChange">
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </label>
        <button type="button" class="contacts-brevo__pager-btn" :disabled="queryParams.page <= 1" @click="goToPage(queryParams.page - 1)">
          Previous
        </button>
        <button type="button" class="contacts-brevo__pager-btn" :disabled="queryParams.page >= totalPages" @click="goToPage(queryParams.page + 1)">
          Next
        </button>
      </div>
    </div>

    <b-modal scroll="keep" :aria-modal="true" :active.sync="isBulkListFormVisible" :width="500" class="has-overflow">
      <subscriber-bulk-list :num-subscribers="this.numSelectedSubscribers" @finished="bulkChangeLists" />
    </b-modal>

    <!-- Edit contact keeps the full modal -->
    <b-modal
      v-if="isEditing"
      scroll="keep"
      :aria-modal="true"
      :active.sync="isFormVisible"
      :width="850"
      @close="onFormClose"
    >
      <subscriber-form :data="curItem" :is-editing="true" @finished="querySubscribers" />
    </b-modal>

    <!-- Create contact — Brevo side drawer -->
    <div
      class="contacts-drawer"
      :class="{ 'is-open': isCreateDrawerOpen }"
      :aria-hidden="isCreateDrawerOpen ? 'false' : 'true'"
    >
      <button type="button" class="contacts-drawer__backdrop" aria-label="Close" @click="closeCreateDrawer" />
      <aside class="contacts-drawer__panel" role="dialog" aria-modal="true" aria-label="Create a contact">
        <header class="contacts-drawer__head">
          <h2>Create a contact</h2>
          <button type="button" class="contacts-drawer__close" aria-label="Close" @click="closeCreateDrawer">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
              <path d="M3 3l8 8M11 3L3 11" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </header>
        <form class="contacts-drawer__body" @submit.prevent="createContactFromDrawer">
          <label class="contacts-drawer__field">
            <span>FIRSTNAME</span>
            <input v-model="createForm.firstname" type="text" placeholder="Enter the FIRSTNAME" maxlength="200" />
          </label>
          <label class="contacts-drawer__field">
            <span>LASTNAME</span>
            <input v-model="createForm.lastname" type="text" placeholder="Enter the LASTNAME" maxlength="200" />
          </label>
          <label class="contacts-drawer__field">
            <span>EMAIL</span>
            <input
              v-model="createForm.email"
              type="email"
              placeholder="Enter the email address"
              maxlength="200"
              required
            />
          </label>
          <label class="contacts-drawer__field">
            <span>SMS</span>
            <input v-model="createForm.sms" type="text" placeholder="+91" maxlength="40" />
          </label>
          <label class="contacts-drawer__field">
            <span>WHATSAPP</span>
            <input v-model="createForm.whatsapp" type="text" placeholder="+91" maxlength="40" />
          </label>
          <label class="contacts-drawer__field">
            <span>LANDLINE_NUMBER</span>
            <input v-model="createForm.landline" type="text" placeholder="+91" maxlength="40" />
          </label>
          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeCreateDrawer">Cancel</button>
            <button type="submit" class="contacts-drawer__save" :disabled="isCreating">Create</button>
          </footer>
        </form>
      </aside>
    </div>

    <!-- Customize columns — Brevo side drawer -->
    <div
      class="contacts-drawer"
      :class="{ 'is-open': isColumnsDrawerOpen }"
      :aria-hidden="isColumnsDrawerOpen ? 'false' : 'true'"
    >
      <button type="button" class="contacts-drawer__backdrop" aria-label="Close" @click="closeColumnsDrawer" />
      <aside class="contacts-drawer__panel" role="dialog" aria-modal="true" aria-label="Attributes visible as columns">
        <header class="contacts-drawer__head">
          <h2>Attributes visible as columns</h2>
          <button type="button" class="contacts-drawer__close" aria-label="Close" @click="closeColumnsDrawer">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
              <path d="M3 3l8 8M11 3L3 11" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </header>
        <div class="contacts-drawer__body contacts-drawer__body--columns">
          <p class="contacts-drawer__help">
            Customize the Contact page, and choose the attributes you want to see as columns.
          </p>

          <div class="contacts-attr-picker" ref="attrPicker">
            <button type="button" class="contacts-attr-picker__trigger" @click.stop="toggleAttrPicker">
              <span class="contacts-attr-picker__plus" aria-hidden="true">+</span>
              Select attributes
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </button>
            <div v-if="isAttrPickerOpen" class="contacts-attr-picker__menu">
              <label class="contacts-attr-picker__search">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <circle cx="11" cy="11" r="6.5" stroke="currentColor" stroke-width="1.8" />
                  <path d="M16 16l4 4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
                </svg>
                <input v-model="attrSearch" type="search" placeholder="Search" aria-label="Search attributes" />
              </label>
              <ul class="contacts-attr-picker__list">
                <li v-for="col in addableColumns" :key="col.key">
                  <button type="button" class="contacts-attr-picker__item" @click="addColumnDraft(col.key)">
                    <span class="contacts-attr-picker__add-icon" aria-hidden="true">+</span>
                    {{ col.label }}
                  </button>
                </li>
                <li v-if="!addableColumns.length" class="contacts-attr-picker__empty">No attributes found</li>
              </ul>
              <button type="button" class="contacts-attr-picker__edit" @click="onEditAttributes">
                Edit &amp; create attributes
                <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                  <path d="M4.5 2.5H2.75A1.25 1.25 0 0 0 1.5 3.75v5.5A1.25 1.25 0 0 0 2.75 10.5h5.5A1.25 1.25 0 0 0 9.5 9.25V7.5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" />
                  <path d="M6.5 1.5H10.5V5.5M10.5 1.5L5.5 6.5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>
          </div>

          <ul class="contacts-columns-list">
            <li
              v-for="(col, idx) in columnOptions"
              :key="col.key"
              class="contacts-columns-list__item"
              :class="{ 'is-dragging': dragIndex === idx }"
              draggable="true"
              @dragstart="onColumnDragStart(idx, $event)"
              @dragover.prevent
              @drop="onColumnDrop(idx)"
            >
              <span class="contacts-columns-list__drag" aria-hidden="true">
                <svg width="10" height="16" viewBox="0 0 10 16" fill="none">
                  <circle cx="3" cy="2" r="1.15" fill="currentColor" />
                  <circle cx="7" cy="2" r="1.15" fill="currentColor" />
                  <circle cx="3" cy="8" r="1.15" fill="currentColor" />
                  <circle cx="7" cy="8" r="1.15" fill="currentColor" />
                  <circle cx="3" cy="14" r="1.15" fill="currentColor" />
                  <circle cx="7" cy="14" r="1.15" fill="currentColor" />
                </svg>
              </span>
              <span class="contacts-columns-list__label">{{ col.label }}</span>
              <button
                type="button"
                class="contacts-columns-list__remove"
                :aria-label="'Remove ' + col.label"
                @click="toggleColumnDraft(col.key)"
              >
                ×
              </button>
            </li>
          </ul>

          <button type="button" class="contacts-columns-edit" @click="onEditAttributes">
            Edit &amp; create attributes
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path d="M4.5 2.5H2.75A1.25 1.25 0 0 0 1.5 3.75v5.5A1.25 1.25 0 0 0 2.75 10.5h5.5A1.25 1.25 0 0 0 9.5 9.25V7.5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" />
              <path d="M6.5 1.5H10.5V5.5M10.5 1.5L5.5 6.5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>

          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeColumnsDrawer">Cancel</button>
            <button type="button" class="contacts-drawer__save contacts-drawer__save--pill" @click="saveColumns">Save</button>
          </footer>
        </div>
      </aside>
    </div>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import dayjs from 'dayjs';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import { uris } from '../constants';
import SubscriberBulkList from './SubscriberBulkList.vue';
import SubscriberForm from './SubscriberForm.vue';

export default Vue.extend({
  components: {
    SubscriberForm,
    SubscriberBulkList,
    EmptyPlaceholder,
  },

  data() {
    return {
      // Current subscriber item being edited.
      curItem: null,
      isSearchAdvanced: false,
      isFilterOpen: false,
      isEditing: false,
      isFormVisible: false,
      isBulkListFormVisible: false,
      isCreateDrawerOpen: false,
      isColumnsDrawerOpen: false,
      isAttrPickerOpen: false,
      attrSearch: '',
      dragIndex: null,
      isCreating: false,
      filterSeq: 1,
      filters: [],

      createForm: {
        firstname: '',
        lastname: '',
        email: '',
        sms: '',
        whatsapp: '',
        landline: '',
      },

      allColumnDefs: [
        { key: 'subscribed', label: 'SUBSCRIBED', sortField: 'status' },
        { key: 'blocklisted', label: 'BLOCKLISTED' },
        { key: 'email', label: 'EMAIL', sortField: 'email' },
        { key: 'firstname', label: 'FIRSTNAME' },
        { key: 'lastname', label: 'LASTNAME' },
        { key: 'sms', label: 'SMS' },
        { key: 'whatsapp', label: 'WHATSAPP' },
        { key: 'landline', label: 'LANDLINE_NUMBER' },
        { key: 'extId', label: 'EXT_ID' },
        { key: 'timezone', label: 'CONTACT_TIMEZONE' },
        { key: 'jobTitle', label: 'JOB_TITLE' },
        { key: 'company', label: 'COMPANY' },
        { key: 'lastChanged', label: 'LAST CHANGED', sortField: 'updated_at' },
        { key: 'created', label: 'CREATION DATE', sortField: 'created_at' },
      ],
      filterFieldDefs: [
        { key: 'email', label: 'Email', type: 'text' },
        { key: 'firstname', label: 'First name', type: 'text' },
        { key: 'lastname', label: 'Last name', type: 'text' },
        { key: 'sms', label: 'SMS', type: 'text' },
        { key: 'whatsapp', label: 'WhatsApp', type: 'text' },
        { key: 'landline', label: 'Landline number', type: 'text' },
        { key: 'ext_id', label: 'EXT_ID', type: 'text' },
        { key: 'timezone', label: 'Contact timezone', type: 'text' },
        { key: 'job_title', label: 'Job title', type: 'text' },
        { key: 'company', label: 'Company', type: 'text' },
        { key: 'list', label: 'Member of a list', type: 'list' },
        { key: 'created', label: 'Creation date', type: 'date' },
        { key: 'updated', label: 'Edit date', type: 'date' },
        { key: 'status', label: 'Email campaigns subscriptions', type: 'status' },
      ],
      visibleColumns: ['subscribed', 'email', 'landline', 'lastChanged', 'created'],
      draftColumns: [],

      // Table bulk row selection states.
      bulk: {
        checked: [],
        all: false,
      },

      queryInput: '',

      // Query params to filter the getSubscribers() API call.
      queryParams: {
        // Search query expression.
        queryExp: '',
        search: '',

        // ID of the list the current subscriber view is filtered by.
        listID: null,
        page: 1,
        perPage: 20,
        orderBy: 'id',
        order: 'desc',
        subStatus: null,
      },
    };
  },

  methods: {
    // Count the lists from which a subscriber has not unsubscribed.
    listCount(lists) {
      return lists.reduce((defVal, item) => (defVal + (item.subscriptionStatus !== 'unsubscribed' ? 1 : 0)), 0);
    },

    contactDisplayName(row) {
      const first = this.$utils.subscriberAttrib(row.attribs, 'firstname');
      const last = this.$utils.subscriberAttrib(row.attribs, 'lastname');
      const name = [first, last].filter(Boolean).join(' ').trim();
      if (name) return name;
      if (row.name) return row.name;
      return (row.email || '').split('@')[0] || 'Contact';
    },

    contactPhone(row) {
      const keys = ['landline_phone', 'landline', 'phone', 'mobile', 'SMS'];
      for (let i = 0; i < keys.length; i += 1) {
        const v = this.$utils.subscriberAttrib(row.attribs, keys[i]);
        if (v) return v;
      }
      return '—';
    },

    formatBrevoDate(raw) {
      if (!raw) return '—';
      const d = dayjs(raw);
      if (!d.isValid()) return '—';
      return d.format('DD/MM/YYYY');
    },

    toggleFilterPanel() {
      this.isFilterOpen = !this.isFilterOpen;
      if (this.isFilterOpen && this.filters.length === 0) {
        this.addFilter('and');
      }
    },

    newFilter(join) {
      this.filterSeq += 1;
      return {
        id: this.filterSeq,
        field: 'email',
        op: 'contains',
        value: '',
        join: join || 'and',
      };
    },

    addFilter(join) {
      this.filters = [...this.filters, this.newFilter(join)];
      this.isFilterOpen = true;
    },

    removeFilter(idx) {
      const next = this.filters.slice();
      next.splice(idx, 1);
      this.filters = next;
      if (!this.filters.length) {
        this.isFilterOpen = false;
        this.clearFilters();
      }
    },

    clearFilters() {
      this.filters = [];
      this.queryParams.queryExp = '';
      this.queryParams.page = 1;
      this.isFilterOpen = false;
      this.querySubscribers();
    },

    filterType(f) {
      const def = this.filterFieldDefs.find((d) => d.key === f.field);
      return (def && def.type) || 'text';
    },

    fieldLabel(key) {
      const def = this.filterFieldDefs.find((d) => d.key === key);
      return (def && def.label) || key;
    },

    isEmptyOp(op) {
      return op === 'is_empty' || op === 'is_not_empty';
    },

    operatorsFor(f) {
      const type = this.filterType(f);
      if (type === 'list') {
        return [
          { value: 'is_member', label: 'is member of' },
          { value: 'is_not_member', label: 'is not member of' },
        ];
      }
      if (type === 'status') {
        return [
          { value: 'is', label: 'is' },
          { value: 'is_not', label: 'is not' },
        ];
      }
      if (type === 'date') {
        return [
          { value: 'on', label: 'on' },
          { value: 'before', label: 'before' },
          { value: 'after', label: 'after' },
          { value: 'last_n_days', label: 'in the last # days' },
          { value: 'is_empty', label: 'is empty' },
          { value: 'is_not_empty', label: 'is not empty' },
        ];
      }
      return [
        { value: 'contains', label: 'contains' },
        { value: 'not_contains', label: 'does not contain' },
        { value: 'is', label: 'is' },
        { value: 'is_not', label: 'is not' },
        { value: 'starts_with', label: 'starts with' },
        { value: 'ends_with', label: 'ends with' },
        { value: 'is_empty', label: 'is empty' },
        { value: 'is_not_empty', label: 'is not empty' },
      ];
    },

    onFilterFieldChange(f) {
      const idx = this.filters.indexOf(f);
      if (idx < 0) return;
      const ops = this.operatorsFor(f);
      const next = {
        ...f,
        op: ops[0].value,
        value: this.filterType(f) === 'status' ? 'subscribed' : '',
      };
      this.$set(this.filters, idx, next);
    },

    escapeSql(v) {
      return String(v || '').replace(/'/g, "''");
    },

    attribCol(key) {
      return `subscribers.attribs->>'${key}'`;
    },

    textSql(col, op, value) {
      const v = this.escapeSql(value);
      if (op === 'is') return `${col} ILIKE '${v}'`;
      if (op === 'is_not') return `(${col} IS NULL OR ${col} NOT ILIKE '${v}')`;
      if (op === 'contains') return `${col} ILIKE '%${v}%'`;
      if (op === 'not_contains') return `(${col} IS NULL OR ${col} NOT ILIKE '%${v}%')`;
      if (op === 'starts_with') return `${col} ILIKE '${v}%'`;
      if (op === 'ends_with') return `${col} ILIKE '%${v}'`;
      if (op === 'is_empty') return `(${col} IS NULL OR ${col} = '')`;
      if (op === 'is_not_empty') return `(${col} IS NOT NULL AND ${col} <> '')`;
      return 'TRUE';
    },

    dateSql(col, op, value) {
      const v = this.escapeSql(value);
      if (op === 'is_empty') return `${col} IS NULL`;
      if (op === 'is_not_empty') return `${col} IS NOT NULL`;
      if (op === 'last_n_days') {
        const n = parseInt(value, 10);
        if (!n || n < 1) return 'TRUE';
        return `${col} >= NOW() - INTERVAL '${n} days'`;
      }
      if (!v) return 'TRUE';
      if (op === 'on') return `${col}::date = '${v}'::date`;
      if (op === 'before') return `${col}::date < '${v}'::date`;
      if (op === 'after') return `${col}::date > '${v}'::date`;
      return 'TRUE';
    },

    compileFilter(f) {
      const type = this.filterType(f);
      if (type === 'list') {
        const id = parseInt(f.value, 10);
        if (!id) return null;
        const exists = `EXISTS (SELECT 1 FROM subscriber_lists sl WHERE sl.subscriber_id = subscribers.id AND sl.list_id = ${id})`;
        return f.op === 'is_not_member' ? `NOT ${exists}` : exists;
      }
      if (type === 'status') {
        const status = f.value === 'blocklisted' ? 'blocklisted' : 'enabled';
        if (f.op === 'is_not') return `subscribers.status <> '${status}'`;
        return `subscribers.status = '${status}'`;
      }
      if (type === 'date') {
        const col = f.field === 'updated' ? 'subscribers.updated_at' : 'subscribers.created_at';
        return this.dateSql(col, f.op, f.value);
      }
      const attribMap = {
        firstname: 'firstname',
        lastname: 'lastname',
        sms: 'SMS',
        whatsapp: 'whatsapp',
        landline: 'landline_phone',
        ext_id: 'EXT_ID',
        timezone: 'CONTACT_TIMEZONE',
        job_title: 'job_title',
        company: 'company',
      };
      let col = 'subscribers.email';
      if (f.field === 'email') col = 'subscribers.email';
      else if (attribMap[f.field]) col = this.attribCol(attribMap[f.field]);
      if (f.field === 'firstname' && f.op !== 'is_empty' && f.op !== 'is_not_empty') {
        const a = this.textSql(this.attribCol('firstname'), f.op, f.value);
        const n = this.textSql('subscribers.name', f.op, f.value);
        return `(${a} OR ${n})`;
      }
      return this.textSql(col, f.op, f.value);
    },

    applyFilters() {
      const parts = [];
      this.filters.forEach((f) => {
        const sql = this.compileFilter(f);
        if (!sql || sql === 'TRUE') return;
        if (!parts.length) {
          parts.push(`(${sql})`);
          return;
        }
        const join = f.join === 'or' ? 'OR' : 'AND';
        parts.push(`${join} (${sql})`);
      });
      this.queryParams.queryExp = parts.join(' ');
      this.queryParams.search = '';
      this.queryParams.page = 1;
      if (this.queryParams.queryExp && !this.$can('subscribers:sql_query')) {
        this.$utils.toast('Your role needs the subscribers query permission to use filters.', 'is-danger');
        return;
      }
      this.querySubscribers();
    },

    columnDisplay(row, key) {
      if (key === 'email') return row.email || '—';
      if (key === 'firstname') {
        return this.$utils.subscriberAttrib(row.attribs, 'firstname') || (row.name || '').split(' ')[0] || '—';
      }
      if (key === 'lastname') {
        const last = this.$utils.subscriberAttrib(row.attribs, 'lastname');
        if (last) return last;
        const parts = String(row.name || '').trim().split(/\s+/);
        return parts.length > 1 ? parts.slice(1).join(' ') : '—';
      }
      if (key === 'sms') return this.$utils.subscriberAttrib(row.attribs, 'SMS') || this.$utils.subscriberAttrib(row.attribs, 'sms') || '—';
      if (key === 'whatsapp') return this.$utils.subscriberAttrib(row.attribs, 'whatsapp') || '—';
      if (key === 'landline') return this.contactPhone(row);
      if (key === 'extId') return this.$utils.subscriberAttrib(row.attribs, 'EXT_ID') || this.$utils.subscriberAttrib(row.attribs, 'ext_id') || '—';
      if (key === 'timezone') return this.$utils.subscriberAttrib(row.attribs, 'CONTACT_TIMEZONE') || this.$utils.subscriberAttrib(row.attribs, 'timezone') || '—';
      if (key === 'jobTitle') return this.$utils.subscriberAttrib(row.attribs, 'job_title') || this.$utils.subscriberAttrib(row.attribs, 'JOB_TITLE') || '—';
      if (key === 'company') return this.$utils.subscriberAttrib(row.attribs, 'company') || '—';
      if (key === 'lastChanged') return this.formatBrevoDate(row.updatedAt);
      if (key === 'created') return this.formatBrevoDate(row.createdAt);
      return '—';
    },

    goToPage(p) {
      const page = Math.min(Math.max(1, p), this.totalPages);
      this.querySubscribers({ page });
    },

    onPerPageChange(e) {
      const perPage = parseInt(e.target.value, 10) || 20;
      this.querySubscribers({ page: 1, perPage });
    },

    toggleAttrPicker() {
      this.isAttrPickerOpen = !this.isAttrPickerOpen;
      if (this.isAttrPickerOpen) this.attrSearch = '';
    },

    onColumnDragStart(idx, e) {
      this.dragIndex = idx;
      if (e.dataTransfer) {
        e.dataTransfer.effectAllowed = 'move';
        e.dataTransfer.setData('text/plain', String(idx));
      }
    },

    onColumnDrop(idx) {
      if (this.dragIndex === null || this.dragIndex === idx) {
        this.dragIndex = null;
        return;
      }
      const next = [...this.draftColumns];
      const [moved] = next.splice(this.dragIndex, 1);
      next.splice(idx, 0, moved);
      this.draftColumns = next;
      this.dragIndex = null;
    },

    onEditAttributes() {
      this.isAttrPickerOpen = false;
      this.$utils.toast('Add custom fields on a contact profile. They can then be shown as columns.');
    },

    onDocClick(e) {
      if (!this.isAttrPickerOpen) return;
      const root = this.$refs.attrPicker;
      if (root && !root.contains(e.target)) {
        this.isAttrPickerOpen = false;
      }
    },

    toggleAdvancedSearch() {
      this.isSearchAdvanced = !this.isSearchAdvanced;
      this.queryParams.search = '';

      // Toggling to simple search.
      if (!this.isSearchAdvanced) {
        this.queryInput = '';
        this.queryParams.queryExp = '';
        this.queryParams.page = 1;
        this.querySubscribers();
        this.$refs.query.focus();
        return;
      }

      // Toggling to advanced search.
      const q = this.queryInput.replace(/'/, "''").trim();
      if (q) {
        if (this.$utils.validateEmail(q)) {
          this.queryParams.queryExp = `email = '${q.toLowerCase()}'`;
        } else {
          this.queryParams.queryExp = `(name ~* '${q}' OR email ~* '${q.toLowerCase()}')`;
        }
      }

      // Toggling to advanced search.
      this.$nextTick(() => {
        this.$refs.queryExp.focus();
      });
    },

    // Mark all subscribers in the query as selected.
    selectAllSubscribers() {
      this.bulk.all = true;
    },

    onTableCheck() {
      // Disable bulk.all selection if there are no rows checked in the table.
      if (this.bulk.checked.length !== this.subscribers.total) {
        this.bulk.all = false;
      }
    },

    // Show the edit list form.
    showEditForm(sub) {
      this.closeCreateDrawer();
      this.curItem = sub;
      this.isEditing = true;
      this.isFormVisible = true;
    },

    // Show the new contact form (Brevo side drawer).
    showNewForm() {
      let lists = [];
      if (this.currentList) {
        lists = [this.currentList];
      } else if (this.queryParams.listID) {
        const found = this.availableLists.find((l) => l.id === this.queryParams.listID);
        lists = found ? [found] : [{ id: this.queryParams.listID }];
      }
      this.curItem = lists.length ? { lists } : {};
      this.isEditing = false;
      this.isFormVisible = false;
      this.createForm = {
        firstname: '',
        lastname: '',
        email: '',
        sms: '',
        whatsapp: '',
        landline: '',
      };
      this.isCreateDrawerOpen = true;
    },

    closeCreateDrawer() {
      this.isCreateDrawerOpen = false;
    },

    createContactFromDrawer() {
      const email = (this.createForm.email || '').trim();
      if (!email) {
        return;
      }
      const first = (this.createForm.firstname || '').trim();
      const last = (this.createForm.lastname || '').trim();
      const name = [first, last].filter(Boolean).join(' ') || email.split('@')[0];
      const attribs = this.$utils.mergeSubscriberAttribs({}, {
        firstname: first,
        lastname: last,
        company: '',
      });
      if (this.createForm.sms) attribs.SMS = this.createForm.sms.trim();
      if (this.createForm.whatsapp) attribs.whatsapp = this.createForm.whatsapp.trim();
      if (this.createForm.landline) attribs.landline_phone = this.createForm.landline.trim();

      const lists = (this.curItem && this.curItem.lists)
        ? this.curItem.lists.map((l) => l.id)
        : [];

      this.isCreating = true;
      this.$api.createSubscriber({
        email,
        name,
        status: 'enabled',
        attribs,
        lists,
        preconfirm_subscriptions: false,
      }).then((d) => {
        this.closeCreateDrawer();
        this.querySubscribers();
        this.$utils.toast(this.$t('globals.messages.created', { name: d.name || email }));
      }).finally(() => {
        this.isCreating = false;
      });
    },

    openColumnsDrawer() {
      this.draftColumns = [...this.visibleColumns];
      this.attrSearch = '';
      this.isAttrPickerOpen = false;
      this.isColumnsDrawerOpen = true;
    },

    closeColumnsDrawer() {
      this.isColumnsDrawerOpen = false;
      this.isAttrPickerOpen = false;
    },

    isColumnVisible(key) {
      return this.visibleColumns.includes(key);
    },

    toggleColumnDraft(key) {
      this.draftColumns = this.draftColumns.filter((k) => k !== key);
    },

    addColumnDraft(key) {
      if (!this.draftColumns.includes(key)) {
        this.draftColumns = [...this.draftColumns, key];
      }
    },

    saveColumns() {
      this.visibleColumns = this.draftColumns.length
        ? [...this.draftColumns]
        : ['email'];
      try {
        localStorage.setItem('nexuses.contacts.columns', JSON.stringify(this.visibleColumns));
      } catch (e) {
        // ignore
      }
      this.closeColumnsDrawer();
    },

    loadSavedColumns() {
      try {
        const raw = localStorage.getItem('nexuses.contacts.columns');
        if (!raw) return;
        const parsed = JSON.parse(raw);
        if (Array.isArray(parsed) && parsed.length) {
          this.visibleColumns = parsed;
        }
      } catch (e) {
        // ignore
      }
    },

    openAddIfRequested() {
      if (this.$route.query.add !== '1' || !this.$can('subscribers:manage')) {
        return;
      }
      this.$nextTick(() => {
        this.showNewForm();
        const q = { ...this.$route.query };
        delete q.add;
        this.$router.replace({ path: this.$route.path, query: q });
      });
    },

    showBulkListForm() {
      this.isBulkListFormVisible = true;
    },

    onFormClose() {
      this.isFormVisible = false;
      if (this.$route.params.id) {
        this.$router.push({ name: 'subscribers' });
      }
    },

    onPageChange(p) {
      this.goToPage(p);
    },

    onSort(field, direction) {
      this.querySubscribers({ orderBy: field, order: direction });
    },

    // Prepares an SQL expression for simple name search inputs and saves it
    // in this.queryExp.
    onSimpleQueryInput(v) {
      const q = v.replace(/'/, "''").trim();
      this.queryParams.page = 1;
      this.queryParams.search = q.toLowerCase();
    },

    // Ctrl + Enter on the advanced query searches.
    onAdvancedQueryEnter(e) {
      if (e.ctrlKey) {
        this.onSubmit();
      }
    },

    onSubmit() {
      this.querySubscribers({ page: 1 });
    },

    onFilterListSelect(listID) {
      const normalizedID = listID ? parseInt(listID, 10) : null;
      this.queryParams.listID = normalizedID;
      this.querySubscribers({ page: 1 });
      if (normalizedID) {
        this.$router.replace({ name: 'subscribers_list', params: { listID: normalizedID } });
      } else {
        this.$router.replace({ name: 'subscribers' });
      }
    },

    fetchFilterLists() {
      this.$api.getLists({ minimal: true, per_page: 'all', status: 'active' }).catch(() => {});
    },

    // Search / query subscribers.
    querySubscribers(params) {
      this.queryParams = { ...this.queryParams, ...params };

      const qp = {
        list_id: this.queryParams.listID,
        search: this.queryParams.search,
        query: this.queryParams.queryExp,
        page: this.queryParams.page,
        per_page: this.queryParams.perPage,
        subscription_status: this.queryParams.subStatus,
        order_by: this.queryParams.orderBy,
        order: this.queryParams.order,
      };

      if (!this.queryParams.queryExp) {
        delete qp.query;
      }
      if (!this.queryParams.search) {
        delete qp.search;
      }

      this.$nextTick(() => {
        this.$api.getSubscribers(qp).then(() => {
          this.bulk.checked = [];
        });
      });
    },

    deleteSubscriber(sub) {
      this.$utils.confirm(
        null,
        () => {
          this.$api.deleteSubscriber(sub.id).then(() => {
            this.querySubscribers();

            this.$utils.toast(this.$t('globals.messages.deleted', { name: sub.name }));
          });
        },
      );
    },

    blocklistSubscribers() {
      let fn = null;
      if (!this.bulk.all && this.bulk.checked.length > 0) {
        // If 'all' is not selected, blocklist subscribers by IDs.
        fn = () => {
          const ids = this.bulk.checked.map((s) => s.id);
          this.$api.blocklistSubscribers({ ids })
            .then(() => this.querySubscribers());
        };
      } else {
        // 'All' is selected, blocklist by query.
        fn = () => {
          this.$api.blocklistSubscribersByQuery({
            search: this.queryParams.search,
            query: this.queryParams.queryExp,
            list_ids: this.queryParams.listID ? [this.queryParams.listID] : null,
            subscription_status: this.queryParams.subStatus,
          }).then(() => this.querySubscribers());
        };
      }

      this.$utils.confirm(this.$t('subscribers.confirmBlocklist', { num: this.numSelectedSubscribers }), fn);
    },

    exportSubscribers() {
      const num = !this.bulk.all && this.bulk.checked.length > 0
        ? this.bulk.checked.length : this.subscribers.total;

      this.$utils.confirm(this.$t('subscribers.confirmExport', { num }), () => {
        const q = new URLSearchParams();

        if (this.queryParams.search) {
          q.append('search', this.queryParams.search);
        } else if (this.queryParams.queryExp) {
          q.append('query', this.queryParams.queryExp);
        }

        if (this.queryParams.listID) {
          q.append('list_id', this.queryParams.listID);
        }

        if (this.queryParams.subStatus) {
          q.append('subscription_status', this.queryParams.subStatus);
        }

        // Export selected subscribers.
        if (!this.bulk.all && this.bulk.checked.length > 0) {
          this.bulk.checked.map((s) => q.append('id', s.id));
        }

        document.location.href = `${uris.exportSubscribers}?${q.toString()}`;
      });
    },

    deleteSubscribers() {
      let fn = null;
      if (!this.bulk.all && this.bulk.checked.length > 0) {
        // If 'all' is not selected, delete subscribers by IDs.
        fn = () => {
          const ids = this.bulk.checked.map((s) => s.id);
          this.$api.deleteSubscribers({ id: ids })
            .then(() => {
              this.querySubscribers();

              this.$utils.toast(this.$t('subscribers.subscribersDeleted', { num: this.numSelectedSubscribers }));
            });
        };
      } else {
        // 'All' is selected, delete by query.
        fn = () => {
          this.$api.deleteSubscribersByQuery({
            // If the query expression is empty, explicitly pass `all=true`
            // so that the backend deletes all records in the DB with an empty query string.
            all: this.queryParams.queryExp.trim() === '' && this.queryParams.search.trim() === '',
            search: this.queryParams.search,
            query: this.queryParams.queryExp,
            list_ids: this.queryParams.listID ? [this.queryParams.listID] : null,
            subscription_status: this.queryParams.subStatus,
          }).then(() => {
            this.querySubscribers();

            this.$utils.toast(this.$t(
              'subscribers.subscribersDeleted',
              { num: this.numSelectedSubscribers },
            ));
          });
        };
      }

      this.$utils.confirm(this.$t('subscribers.confirmDelete', { num: this.numSelectedSubscribers }), fn);
    },

    bulkChangeLists(action, preconfirm, lists) {
      const data = {
        action,
        query: this.fullQueryExp,
        search: this.queryParams.search,
        list_ids: this.queryParams.listID ? [this.queryParams.listID] : null,
        target_list_ids: lists.map((l) => l.id),
      };

      if (preconfirm) {
        data.status = 'confirmed';
      }

      let fn = null;
      if (!this.bulk.all && this.bulk.checked.length > 0) {
        // If 'all' is not selected, perform by IDs.
        fn = this.$api.addSubscribersToLists;
        data.ids = this.bulk.checked.map((s) => s.id);
      } else {
        // 'All' is selected, perform by query.
        data.query = this.queryParams.queryExp;
        data.subscription_status = this.queryParams.subStatus;
        fn = this.$api.addSubscribersToListsByQuery;
      }

      fn(data).then(() => {
        this.querySubscribers();
        this.$utils.toast(this.$t('subscribers.listChangeApplied'));
      });
    },
  },

  computed: {
    ...mapState(['subscribers', 'lists', 'loading']),

    numSelectedSubscribers() {
      if (this.bulk.all) {
        return this.subscribers.total;
      }
      return this.bulk.checked.length;
    },

    // Returns the list that the subscribers are being filtered by in.
    currentList() {
      if (!this.queryParams.listID || !this.lists.results) {
        return null;
      }

      return this.lists.results.find((l) => l.id === this.queryParams.listID);
    },

    availableLists() {
      return (this.lists && this.lists.results) || [];
    },

    fullQueryExp() {
      return this.queryParams.queryExp || '';
    },

    importRoute() {
      if (this.currentList) {
        return { name: 'import', query: { list_id: this.currentList.id } };
      }
      return { name: 'import' };
    },

    columnOptions() {
      return this.draftColumns
        .map((key) => this.allColumnDefs.find((c) => c.key === key))
        .filter(Boolean);
    },

    visibleColumnDefs() {
      return this.visibleColumns
        .map((key) => this.allColumnDefs.find((c) => c.key === key))
        .filter(Boolean);
    },

    addableColumns() {
      const q = (this.attrSearch || '').trim().toLowerCase();
      return this.allColumnDefs.filter((c) => {
        if (this.draftColumns.includes(c.key)) return false;
        if (!q) return true;
        return c.label.toLowerCase().indexOf(q) > -1 || c.key.toLowerCase().indexOf(q) > -1;
      });
    },

    totalPages() {
      const per = this.subscribers.perPage || this.queryParams.perPage || 20;
      const total = this.subscribers.total || 0;
      return Math.max(1, Math.ceil(total / per));
    },

    pagerSummary() {
      const total = this.subscribers.total || 0;
      const per = this.subscribers.perPage || this.queryParams.perPage || 20;
      if (!total) return '0 contacts';
      const start = ((this.queryParams.page - 1) * per) + 1;
      const end = Math.min(this.queryParams.page * per, total);
      return `Showing ${start}–${end} of ${this.$utils.formatNumber(total)} contacts`;
    },
  },

  watch: {
    $route(to, from) {
      if (to.name !== 'subscribers' && to.name !== 'subscribers_list' && to.name !== 'subscriber') {
        return;
      }
      const listID = to.params.listID ? parseInt(to.params.listID, 10) : null;
      const prevListID = from && from.params.listID ? parseInt(from.params.listID, 10) : null;
      if (listID !== prevListID || listID !== this.queryParams.listID) {
        this.queryParams.listID = listID;
        if (!to.params.id) {
          this.querySubscribers({ page: 1 });
        }
      }
      this.openAddIfRequested();
    },
  },

  created() {
    this.$root.$on('page.refresh', this.querySubscribers);
  },

  destroyed() {
    this.$root.$off('page.refresh', this.querySubscribers);
    document.removeEventListener('click', this.onDocClick);
  },

  mounted() {
    document.addEventListener('click', this.onDocClick);
    this.loadSavedColumns();
    this.fetchFilterLists();

    if (this.$route.params.listID) {
      this.queryParams.listID = parseInt(this.$route.params.listID, 10);
    }
    if (this.$route.query.subscription_status) {
      this.queryParams.subStatus = this.$route.query.subscription_status;
    }

    if (this.$route.params.id) {
      this.$api.getSubscriber(parseInt(this.$route.params.id, 10)).then((data) => {
        this.showEditForm(data);
      });
    } else {
      this.querySubscribers();
      this.openAddIfRequested();
    }
  },
});
</script>

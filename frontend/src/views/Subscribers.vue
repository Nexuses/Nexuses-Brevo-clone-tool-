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

        <button type="button" class="contacts-brevo__chip" @click.prevent="toggleAdvancedSearch">
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
              :disabled="isSearchAdvanced"
              data-cy="search"
            />
          </b-field>
        </form>
      </div>
    </div>

    <div v-if="isSearchAdvanced" class="contacts-brevo__advanced mb-4">
      <b-input v-model="queryParams.queryExp" @keydown.native.enter="onAdvancedQueryEnter" type="textarea"
        ref="queryExp" placeholder="subscribers.name LIKE '%user%' or subscribers.status='blocklisted'"
        data-cy="query" />
      <div class="buttons mt-2">
        <b-button type="is-primary" icon-left="magnify" @click="onSubmit" data-cy="btn-query">
          {{ $t('subscribers.query') }}
        </b-button>
        <b-button @click.prevent="toggleAdvancedSearch" icon-left="cancel" data-cy="btn-query-reset">
          {{ $t('subscribers.reset') }}
        </b-button>
      </div>
    </div>

    <b-table :data="subscribers.results ?? []" :loading="loading.subscribers" @check-all="onTableCheck"
      @check="onTableCheck" :checked-rows.sync="bulk.checked" paginated backend-pagination pagination-position="bottom"
      @page-change="onPageChange" :current-page="queryParams.page" :per-page="subscribers.perPage"
      :total="subscribers.total" hoverable checkable backend-sorting @sort="onSort"
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

      <b-table-column v-if="isColumnVisible('subscribed')" v-slot="props" field="status" label="SUBSCRIBED"
        header-class="cy-status">
        <span
          class="contacts-brevo__sub-pill"
          :class="{ 'is-blocklisted': props.row.status === 'blocklisted' }"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M4 6h16v12H4V6z" stroke="currentColor" stroke-width="1.8" />
            <path d="M4 7l8 6 8-6" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" />
          </svg>
          Email
        </span>
      </b-table-column>

      <b-table-column v-if="isColumnVisible('email')" v-slot="props" field="email" label="EMAIL"
        header-class="cy-email-col" sortable>
        <span class="contacts-brevo__email" :title="props.row.email">{{ props.row.email }}</span>
      </b-table-column>

      <b-table-column v-if="isColumnVisible('landline')" v-slot="props" field="phone" label="LANDLINE_PHONE"
        header-class="cy-phone">
        <span class="contacts-brevo__cell">{{ contactPhone(props.row) }}</span>
      </b-table-column>

      <b-table-column v-if="isColumnVisible('lastChanged')" v-slot="props" field="updated_at" label="LAST CHANGED"
        header-class="cy-updated" sortable>
        <span class="contacts-brevo__cell">{{ formatBrevoDate(props.row.updatedAt) }}</span>
      </b-table-column>

      <b-table-column v-if="isColumnVisible('created')" v-slot="props" field="created_at" label="CREATED"
        header-class="cy-created" sortable>
        <span class="contacts-brevo__cell">{{ formatBrevoDate(props.row.createdAt) }}</span>
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
        <div class="contacts-drawer__body">
          <p class="contacts-drawer__help">
            Customize the Contact page, and choose the attributes you want to see as columns.
          </p>
          <ul class="contacts-columns-list">
            <li v-for="col in columnOptions" :key="col.key" class="contacts-columns-list__item">
              <span class="contacts-columns-list__drag" aria-hidden="true">⋮⋮</span>
              <span class="contacts-columns-list__label">{{ col.label }}</span>
              <button
                type="button"
                class="contacts-columns-list__remove"
                :aria-label="`Remove ${col.label}`"
                @click="toggleColumnDraft(col.key)"
              >
                ×
              </button>
            </li>
          </ul>
          <div class="contacts-columns-add">
            <b-dropdown>
              <template #trigger>
                <button type="button" class="contacts-columns-add__btn">
                  <span aria-hidden="true">+</span> Select attributes
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
              </template>
              <b-dropdown-item
                v-for="col in allColumnDefs"
                :key="col.key"
                @click="addColumnDraft(col.key)"
              >
                <span :class="{ 'has-text-grey': draftColumns.includes(col.key) }">
                  {{ col.label }}
                  <template v-if="draftColumns.includes(col.key)"> ✓</template>
                </span>
              </b-dropdown-item>
            </b-dropdown>
          </div>
          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeColumnsDrawer">Cancel</button>
            <button type="button" class="contacts-drawer__save" @click="saveColumns">Save</button>
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
      isEditing: false,
      isFormVisible: false,
      isBulkListFormVisible: false,
      isCreateDrawerOpen: false,
      isColumnsDrawerOpen: false,
      isCreating: false,

      createForm: {
        firstname: '',
        lastname: '',
        email: '',
        sms: '',
        whatsapp: '',
        landline: '',
      },

      allColumnDefs: [
        { key: 'subscribed', label: 'SUBSCRIBED' },
        { key: 'email', label: 'EMAIL' },
        { key: 'landline', label: 'LANDLINE_PHONE' },
        { key: 'lastChanged', label: 'LAST CHANGED' },
        { key: 'created', label: 'CREATION DATE' },
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
      this.isColumnsDrawerOpen = true;
    },

    closeColumnsDrawer() {
      this.isColumnsDrawerOpen = false;
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
      this.querySubscribers({ page: p });
    },

    onSort(field, direction) {
      this.querySubscribers({ orderBy: field, order: direction });
    },

    // Prepares an SQL expression for simple name search inputs and saves it
    // in this.queryExp.
    onSimpleQueryInput(v) {
      const q = v.replace(/'/, "''").trim();
      this.queryParams.queryExp = '';
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
        subscription_status: this.queryParams.subStatus,
        order_by: this.queryParams.orderBy,
        order: this.queryParams.order,
      };

      if (this.queryParams.queryExp) {
        delete qp.search;
      } else {
        delete qp.query;
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

    importRoute() {
      if (this.currentList) {
        return { name: 'import', query: { list_id: this.currentList.id } };
      }
      return { name: 'import' };
    },

    columnOptions() {
      return this.allColumnDefs.filter((c) => this.draftColumns.includes(c.key));
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
  },

  mounted() {
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

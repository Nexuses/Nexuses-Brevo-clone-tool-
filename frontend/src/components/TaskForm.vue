<template>
  <div class="task-modal">
    <header class="task-modal__head">
      <h3>{{ isEditing ? 'Edit task' : 'Create a task' }}</h3>
      <button type="button" class="task-modal__close" aria-label="Close" @click="$emit('close')">×</button>
    </header>

    <div class="task-modal__body">
      <b-field label="Task type">
        <b-select v-model="form.type" expanded>
          <option value="todo">To do</option>
          <option value="call">Call</option>
          <option value="email">Email</option>
          <option value="meeting">Meeting</option>
        </b-select>
      </b-field>

      <b-field label="Name">
        <b-input v-model="form.name" maxlength="200" placeholder="Task name" required />
      </b-field>

      <div class="columns is-mobile">
        <div class="column">
          <b-field label="Due date *">
            <b-datepicker
              v-model="form.dueDate"
              icon="calendar-clock"
              :editable="false"
              placeholder="Due date"
              required
            />
          </b-field>
        </div>
        <div class="column">
          <b-field label="Time">
            <b-select v-model="form.time" expanded>
              <option v-for="t in timeOptions" :key="t" :value="t">{{ t }}</option>
            </b-select>
          </b-field>
        </div>
      </div>

      <div class="task-modal__toggle">
        <span class="task-modal__toggle-lbl">Set reminder</span>
        <label class="task-switch">
          <input type="checkbox" v-model="form.reminder" />
          <span class="task-switch__slider" />
        </label>
        <span class="task-modal__toggle-hint">{{ form.reminder ? 'Reminder on' : 'No reminder' }}</span>
      </div>

      <div class="task-modal__toggle">
        <span class="task-modal__toggle-lbl">Task priority</span>
        <label class="task-switch">
          <input type="checkbox" v-model="form.highPriority" />
          <span class="task-switch__slider" />
        </label>
        <span class="task-modal__toggle-hint">Set task to high priority</span>
      </div>

      <b-field label="Notes">
        <b-input
          v-model="form.notes"
          type="textarea"
          maxlength="10000"
          placeholder="Take notes here..."
          rows="5"
        />
        <p class="help has-text-right">{{ (form.notes || '').length }}/10000</p>
      </b-field>

      <a href="#" class="task-modal__associate" @click.prevent>
        Associate task ({{ form.associations }})
      </a>
    </div>

    <footer class="task-modal__foot">
      <button type="button" class="task-modal__cancel" @click="$emit('close')">Cancel</button>
      <button type="button" class="task-modal__create" @click="submit">
        {{ isEditing ? 'Save' : 'Create' }}
      </button>
    </footer>
  </div>
</template>

<script>
import dayjs from 'dayjs';

const STORAGE_KEY = 'nexuses.crm.tasks';

export function loadTasks() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch (e) {
    return [];
  }
}

export function saveTasks(tasks) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(tasks));
}

export default {
  name: 'TaskForm',

  props: {
    data: { type: Object, default: null },
    defaultDate: { type: Date, default: null },
  },

  data() {
    let due = new Date();
    if (this.data && this.data.dueDate) {
      due = new Date(this.data.dueDate);
    } else if (this.defaultDate) {
      due = new Date(this.defaultDate);
    }
    const base = this.data || {};

    return {
      form: {
        id: base.id || null,
        type: base.type || 'todo',
        name: base.name || 'To do',
        dueDate: due,
        time: base.time || '4:00 PM',
        reminder: !!base.reminder,
        highPriority: !!base.highPriority,
        notes: base.notes || '',
        associations: base.associations || 0,
        status: base.status || 'open',
        createdAt: base.createdAt || null,
      },
      timeOptions: [
        '8:00 AM', '9:00 AM', '10:00 AM', '11:00 AM', '12:00 PM',
        '1:00 PM', '2:00 PM', '3:00 PM', '4:00 PM', '5:00 PM',
        '6:00 PM', '7:00 PM', '8:00 PM',
      ],
    };
  },

  computed: {
    isEditing() {
      return !!this.form.id;
    },
  },

  methods: {
    submit() {
      const name = (this.form.name || '').trim();
      if (!name) {
        this.$utils.toast('Please enter a task name');
        return;
      }
      if (!this.form.dueDate) {
        this.$utils.toast('Please select a due date');
        return;
      }

      const tasks = loadTasks();
      const payload = {
        id: this.form.id || (tasks.reduce((m, t) => Math.max(m, t.id || 0), 0) + 1),
        type: this.form.type,
        name,
        dueDate: dayjs(this.form.dueDate).format('YYYY-MM-DD'),
        time: this.form.time,
        reminder: this.form.reminder,
        highPriority: this.form.highPriority,
        notes: this.form.notes || '',
        associations: this.form.associations || 0,
        status: this.form.status || 'open',
        createdAt: this.form.createdAt || new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      if (this.isEditing) {
        const next = tasks.map((t) => (t.id === payload.id ? payload : t));
        saveTasks(next);
        this.$utils.toast('Task updated');
      } else {
        tasks.unshift(payload);
        saveTasks(tasks);
        this.$utils.toast('Task created');
      }

      this.$emit('finished', payload);
      this.$emit('close');
    },
  },
};
</script>

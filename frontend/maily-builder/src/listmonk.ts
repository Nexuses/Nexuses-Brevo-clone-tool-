import type { JSONContent } from '@tiptap/core';
import { Maily } from '@maily-to/render';

/** Default empty Maily document. */
export const DEFAULT_CONTENT: JSONContent = {
  type: 'doc',
  content: [
    {
      type: 'paragraph',
      content: [{ type: 'text', text: '' }],
    },
  ],
};

/** Map Maily variable ids to listmonk Go template expressions. */
const LISTMONK_VARIABLES: Record<string, string> = {
  email: '{{ .Subscriber.Email }}',
  name: '{{ .Subscriber.Name }}',
  first_name: '{{ .Subscriber.FirstName }}',
  firstName: '{{ .Subscriber.FirstName }}',
  last_name: '{{ .Subscriber.LastName }}',
  lastName: '{{ .Subscriber.LastName }}',
  company: '{{ .Subscriber.Company }}',
  uuid: '{{ .Subscriber.UUID }}',
  status: '{{ .Subscriber.Status }}',
};

/** Suggested variables shown in the Maily variable picker. */
export const LISTMONK_VARIABLE_SUGGESTIONS = [
  { name: 'email', required: false },
  { name: 'name', required: false },
  { name: 'first_name', required: false },
  { name: 'last_name', required: false },
  { name: 'company', required: false },
  { name: 'uuid', required: false },
];

function formatListmonkVariable(variable: string, fallback?: string): string {
  const key = variable.trim();
  const tpl = LISTMONK_VARIABLES[key];
  if (tpl) {
    return tpl;
  }
  // Custom subscriber attributes: {{ .Subscriber.Attribs.city }}
  if (key.startsWith('attribs.') || key.startsWith('Attribs.')) {
    const field = key.replace(/^attribs?\./i, '');
    return `{{ .Subscriber.Attribs.${field} }}`;
  }
  if (/^[a-zA-Z_][\w.]*$/.test(key)) {
    return `{{ .Subscriber.Attribs.${key} }}`;
  }
  if (fallback) {
    return `{{ .Subscriber.Attribs.${key} | default "${fallback}" }}`;
  }
  return `{{ .Subscriber.Attribs.${key} }}`;
}

/** Convert Maily JSON to listmonk-ready HTML (Go templates preserved). */
export async function renderListmonkHtml(content: JSONContent): Promise<string> {
  const maily = new Maily(content);
  maily.setShouldReplaceVariableValues(false);
  maily.setVariableFormatter(({ variable, fallback }) => formatListmonkVariable(variable, fallback));
  return maily.render();
}

import React from 'react';
import ReactDOM from 'react-dom/client';
import type { JSONContent } from '@tiptap/core';

import App, { DEFAULT_CONTENT } from './App';

export interface MailyBuilderProps {
  data?: JSONContent;
  onChange?: (json: JSONContent, html: string) => void;
}

function isRendered(containerId: string): boolean {
  const container = document.getElementById(containerId);
  if (!container) {
    return false;
  }
  return container.hasChildNodes();
}

const roots = new Map<string, ReturnType<typeof ReactDOM.createRoot>>();

function render(containerId: string, props: MailyBuilderProps, force = false) {
  const container = document.getElementById(containerId);
  if (!container) {
    console.error(`MailyBuilder: container #${containerId} not found`);
    return;
  }

  if (force && roots.has(containerId)) {
    roots.get(containerId)?.unmount();
    roots.delete(containerId);
    container.innerHTML = '';
  }

  if (!isRendered(containerId)) {
    const root = ReactDOM.createRoot(container);
    roots.set(containerId, root);
    root.render(
      <React.StrictMode>
        <App {...props} />
      </React.StrictMode>,
    );
  }
}

function setContent(containerId: string, data: JSONContent, onChange?: MailyBuilderProps['onChange']) {
  render(containerId, { data, onChange }, true);
}

export { App, render, setContent, isRendered, DEFAULT_CONTENT };

import { useCallback, useEffect, useRef } from 'react';
import type { Editor as TiptapEditor, JSONContent } from '@tiptap/core';
import { Editor } from '@maily-to/core';
import '@maily-to/core/style.css';

import { DEFAULT_CONTENT, LISTMONK_VARIABLE_SUGGESTIONS, renderListmonkHtml } from './listmonk';

export interface AppProps {
  data?: JSONContent;
  onChange?: (json: JSONContent, html: string) => void;
}

export default function App({ data, onChange }: AppProps) {
  const editorRef = useRef<TiptapEditor | null>(null);
  const onChangeRef = useRef(onChange);
  onChangeRef.current = onChange;

  const emitChange = useCallback(async (editor: TiptapEditor) => {
    const json = editor.getJSON();
    const html = await renderListmonkHtml(json);
    onChangeRef.current?.(json, html);
  }, []);

  const handleEditor = useCallback((editor: TiptapEditor) => {
    editorRef.current = editor;
    emitChange(editor);
  }, [emitChange]);

  useEffect(() => {
    if (!data || !editorRef.current) {
      return;
    }
    const current = JSON.stringify(editorRef.current.getJSON());
    const incoming = JSON.stringify(data);
    if (current !== incoming) {
      editorRef.current.commands.setContent(data, false);
    }
  }, [data]);

  return (
    <div className="maily-builder-root">
      <Editor
        contentJson={data ?? DEFAULT_CONTENT}
        onCreate={handleEditor}
        onUpdate={handleEditor}
        config={{
          hasMenuBar: true,
          spellCheck: false,
          autofocus: false,
          immediatelyRender: true,
        }}
        blocks={[
          {
            title: 'Listmonk',
            commands: [
              {
                title: 'Subscriber email',
                searchTerms: ['email', 'subscriber'],
                command: ({ editor, range }) => {
                  editor.chain().focus().deleteRange(range).insertContent({
                    type: 'variable',
                    attrs: { id: 'email', fallback: null, showIfKey: null },
                  }).run();
                },
              },
              {
                title: 'Subscriber name',
                searchTerms: ['name', 'subscriber'],
                command: ({ editor, range }) => {
                  editor.chain().focus().deleteRange(range).insertContent({
                    type: 'variable',
                    attrs: { id: 'name', fallback: null, showIfKey: null },
                  }).run();
                },
              },
            ],
          },
        ]}
      />
    </div>
  );
}

export { DEFAULT_CONTENT, LISTMONK_VARIABLE_SUGGESTIONS };

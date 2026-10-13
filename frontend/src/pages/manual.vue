<template>
    <div class="manual-layout">
        <aside class="manual-sidebar">
            <MarkdownToc :content="markdownContent" />
        </aside>
        <main class="manual-content">
            <div class="markdown-page markdown-body" v-html="renderedContent"></div>
        </main>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { marked } from 'marked'
import MarkdownToc from '@/components/manual/MarkdownToc.vue'

const markdownContent = ref('')
const renderedContent = ref('')

onMounted(async () => {
    try {
        const response = await fetch('/src/assets/manual/manual.md')
        const text = await response.text()
        markdownContent.value = text

        const processedText = text.replace(/!\[(.*?)\]\((.*?)\)/g, (match, alt, src) => {
            if (!src.startsWith('http') && !src.startsWith('/')) {
                return `![${alt}](/src/assets/manual/${src})`
            }
            return match
        })

        renderedContent.value = marked(processedText, {
            breaks: true,
            gfm: true,
            headerIds: true
        })
    } catch (error) {
        console.error('Error loading markdown:', error)
    }
})
</script>

<style scoped>
.manual-layout {
    display: flex;
    width: 100%;
    max-width: 1800px;
    margin: 0 auto;
    padding: 20px;
    gap: 30px;
    height: 100vh;
}

.manual-sidebar {
    flex: 2;
    position: sticky;
    top: 24px;
    height: calc(100vh - 48px);
    overflow-y: auto;
    padding-right: 16px;
}

.manual-content {
    flex: 7;
    padding: 32px;
    height: 100%;
    overflow-y: auto;
}

.manual-sidebar::-webkit-scrollbar {
    display: none;
}

.manual-content::-webkit-scrollbar {
    display: none;
}

.manual-sidebar::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 2px;
}

:deep(.markdown-body) {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    font-size: 18px;
    line-height: 2.2;
    color: #24292e;
}

:deep(.markdown-body h1),
:deep(.markdown-body h2),
:deep(.markdown-body h3),
:deep(.markdown-body h4),
:deep(.markdown-body h5),
:deep(.markdown-body h6) {
    margin-top: 1.5em;
    margin-bottom: 0.8em;
    font-weight: 600;
    line-height: 1.25;
}

:deep(.markdown-body h1) {
    font-size: 2em;
    border-bottom: 1px solid #eaecef;
    padding-bottom: 0.3em;
}

:deep(.markdown-body h2) {
    font-size: 1.5em;
    border-bottom: 1px solid #eaecef;
    padding-bottom: 0.3em;
}

:deep(.markdown-body code) {
    background-color: rgba(27, 31, 35, 0.05);
    border-radius: 3px;
    font-size: 85%;
    padding: 0.2em 0.4em;
    font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
}

:deep(.markdown-body pre) {
    background-color: #f5f5f5;
    border-radius: 6px;
    padding: 16px;
    overflow: auto;
    line-height: 1.45;
}

:deep(.markdown-body pre code) {
    background-color: #f5f5f5;
    padding: 0;
    border-radius: 0;
}

:deep(.markdown-body blockquote) {
    border-left: 4px solid #1a5fa4;
    background-color: #f5f5f5;
    color: #333;
    padding: 1em 1.5em;
    margin: 1em 0;
    border-radius: 4px;
}

:deep(.markdown-body img) {
    max-width: 100%;
    box-sizing: initial;
    background-color: #fff;
    border-radius: 6px;
    margin: 1.5em 0;
}

:deep(.markdown-body table) {
    border-collapse: collapse;
    width: 100%;
    margin: 1em 0;
    display: block;
    overflow: auto;
}

:deep(.markdown-body table th) {
    font-weight: 600;
    background-color: #f6f8fa;
}

:deep(.markdown-body table th),
:deep(.markdown-body table td) {
    padding: 6px 13px;
    border: 1px solid #dfe2e5;
}

:deep(.markdown-body ul),
:deep(.markdown-body ol) {
    padding-left: 2em;
    margin-bottom: 1em;
}

:deep(.markdown-body li) {
    margin-bottom: 0.5em;
}
</style>

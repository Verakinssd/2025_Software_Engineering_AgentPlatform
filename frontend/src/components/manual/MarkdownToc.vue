<template>
  <div class="toc-wrapper">
    <div class="toc-timeline">
      <div
        v-for="item in headers"
        :key="item.slug"
        :class="['toc-item', { active: activeSlug === item.slug }, `level-${item.level}`]"
      >
        <a
          href="#"
          @click.prevent="scrollTo(item.slug)"
          class="toc-link"
        >
          {{ item.title }}
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'

const props = defineProps({
    content: String
})

const headers = ref([])
const activeSlug = ref('')

const observer = ref(null)

const updateHeaders = () => {
    const markdownPage = document.querySelector('.markdown-page')
    if (!markdownPage) {
        headers.value = []
        return
    }

    const headingElements = markdownPage.querySelectorAll('h1, h2, h3')
    const newHeaders = Array.from(headingElements).map((el, index) => {
        const raw = el.innerText.trim()
        const slug = `section-${index}-${raw.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')}`
        if (!el.id) el.id = slug

        // 添加 level 识别
        const level = parseInt(el.tagName.substring(1))

        return {
            title: raw.replace(/^[#\s]+/, '').trim(),
            slug,
            offsetTop: el.offsetTop,
            level
        }
    })

    headers.value = newHeaders
}

const scrollTo = (slug) => {
    const el = document.getElementById(slug)
    if (el) {
        el.scrollIntoView({ behavior: 'smooth' })
    }
}

const updateActiveSlug = () => {
    if (!headers.value.length) return

    const scrollTop = window.scrollY + 100
    const current = [...headers.value].reverse().find(h => h.offsetTop <= scrollTop)
    if (current) {
        activeSlug.value = current.slug
    }
}

onMounted(() => {
    updateHeaders()

    observer.value = new MutationObserver(updateHeaders)
    const markdownPage = document.querySelector('.markdown-page')
    if (markdownPage) {
        observer.value.observe(markdownPage, {
            childList: true,
            subtree: true
        })
    }

    window.addEventListener('scroll', updateActiveSlug)

    setTimeout(updateHeaders, 500)
})

onBeforeUnmount(() => {
    if (observer.value) observer.value.disconnect()
    window.removeEventListener('scroll', updateActiveSlug)
})

watch(() => props.content, () => {
    setTimeout(updateHeaders, 100)
})
</script>

<style scoped>
.toc-wrapper {
    position: sticky;
    top: 1rem;
    padding: 2rem 0;
}

.toc-wrapper::-webkit-scrollbar {
    display: none;
}

.toc-timeline {
    position: relative;
    padding-left: 1.5rem;
}

.toc-timeline::before {
    content: "";
    position: absolute;
    left: 0.875rem;
    top: 0;
    bottom: 0;
    width: 2px;
    background-color: #e0e0e0;
    z-index: 0;
}

.toc-item {
    position: relative;
    margin-bottom: 1rem;
    padding-left: 0.8rem;
    transition: all 0.2s ease;
}


.toc-item.level-1 { padding-left: 0.5rem; }
.toc-item.level-2 { padding-left: 1.5rem; }
.toc-item.level-3 { padding-left: 3.0rem; }

.toc-item::before {
    content: "";
    position: absolute;
    left: -0.6rem;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: white;
    border: 2px solid #e0e0e0;
    box-sizing: border-box;
    z-index: 1;
    transition: all 0.2s ease;
}

.toc-link {
    display: block;
    color: #546e7a;
    text-decoration: none;
    font-size: 1.1rem;
    line-height: 1.4;
    transition: all 0.2s ease;
}

.toc-item.level-2 .toc-link {
  font-size: 1rem;
}

.toc-item.level-3 .toc-link {
  font-size: 0.9rem;
}

.toc-item:hover::before {
    border-color: #1a5fa4;
    transform: translate(-50%, -50%) scale(1.1);
}

.toc-item:hover .toc-link {
    color: #1a5fa4;
    font-weight: 500;
}

.toc-item.active::before {
    background-color: #1a5fa4;
    border-color: #1a5fa4;
    transform: translate(-50%, -50%) scale(1.1);
    box-shadow: 0 0 0 2px rgba(25, 118, 210, 0.2);
}

.toc-item.active .toc-link {
    color: #1a5fa4;
    font-weight: 500;
}
</style>

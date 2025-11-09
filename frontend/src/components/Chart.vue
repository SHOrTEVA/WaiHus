<script setup lang="ts">
import { onMounted, ref } from 'vue';
import axios from 'axios';

const isBarSelected = ref(true);
const isTableSelected = ref(false);
const characters = ref([]);
const searchParams = defineModels();
onMounted(async () => {
    try {
        const params = {
            q: "Wukong",
            page: 1,
            limit: 5,
            order_by: "favorites",
            sort: "desc"
        };
        const response = await axios.get('https://api.jikan.moe/v4/characters', { searchParams });
        characters.value = response.data.data;
        console.log(characters.value);
    } catch (err) {
        console.error(err);
    }
});

</script>

<template>
    <div class="ChartTab">
        <ul>
            <li 
                :class="{ active: isBarSelected }" 
                @click="isBarSelected = true; isTableSelected = false"
            >
                Bar Chart
            </li>
            <li 
                :class="{ active: isTableSelected }" 
                @click="isBarSelected = false; isTableSelected = true"
            >
                Table
            </li>
        </ul>
    </div>
    <div v-if="isBarSelected" class="bar-chart">Bar Chart Content</div>
    <div v-if="isTableSelected" class="table-chart">Table Content</div>
    <div class="chart-container">
        <h2>Chart Component</h2>
        <form class="search-bar" @submit.prevent="fetchCharacters">
            <input
                v-model="q"
                type="text"
                placeholder="Search characters (e.g. Wukong)"
                aria-label="Search characters"
            />
            <button type="submit">Search</button>
        </form>
        <div v-if="characters && characters.length" class="characters-list">
            <ul>
                <li v-for="char in characters" :key="char.mal_id" class="character-item">
                    <img :src="char.images.jpg.image_url" :alt="char.name" class="avatar" />
                    <div class="info">
                        <h3><a :href="char.url" target="_blank" rel="noopener">{{ char.name }}</a></h3>
                        <p class="favorites">❤️ {{ char.favorites }}</p>
                        <p class="about" v-if="char.about">{{ char.about }}</p>
                    </div>
                </li>
            </ul>
        </div>

        <div v-else class="empty">No characters to display.</div>
    </div>
</template>

<style scoped>
.chart-container {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
}
</style>
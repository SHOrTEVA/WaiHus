<script setup lang="ts">
import { onMounted, ref } from 'vue';
import axios from 'axios';

interface SearchParams {
    q: string;
    page: number;
    limit: number;
    order_by: string;
    sort: string;
}

const isBarSelected = ref(true);
const isTableSelected = ref(false);
const characters = ref<any[]>([]);
const report = ref<any[]>([]);
const searchParams = ref<SearchParams>({
    q: "",
    page: 1,
    limit: 5,
    order_by: "favorites",
    sort: "desc"
});
const host = 'http://localhost:3000'

const fetchCharacters = async () => {
    try {
        const response = await axios.get('https://api.jikan.moe/v4/characters', { params: searchParams.value });
        characters.value = response.data.data;
        console.log(characters.value);
    } catch (err) {
        console.error(err);
    }
};

const submitVote = async (char: any) => {
    try {
        await axios.post(`${host}/api/vote`, {
            character_id: String(char.mal_id),
            name: char.name,
            image_url: char.images?.jpg?.image_url || null
        });
        // Refresh report after voting
        await fetchReport();
        alert(`Thanks for voting for ${char.name}`);
    } catch (err) {
        console.error(err);
        alert('Failed to submit vote');
    }
};

const fetchReport = async () => {
    try {
        const res = await axios.get(`${host}/api/report`);
        report.value = res.data;
    } catch (err) {
        console.error(err);
    }
};

const downloadCSV = () => {
    // trigger backend CSV download
    window.open(`${host}/api/export`, '_blank');
};

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
                v-model="searchParams.q"
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
                        <button @click="submitVote(char)" class="vote-btn">Vote</button>
                    </div>
                </li>
            </ul>
        </div>

        <div v-else class="empty">No characters to display.</div>

        <hr />
        <section class="report-section">
            <h3>Voting Report</h3>
            <button @click="fetchReport">Refresh Report</button>
            <button @click="downloadCSV">Download CSV</button>

            <div v-if="report && report.length" class="report-list">
                <ul>
                    <li v-for="r in report" :key="r.characterId" class="report-item">
                        <div class="r-info">
                            <img v-if="r.imageUrl" :src="r.imageUrl" :alt="r.name" class="avatar-sm" />
                            <strong>{{ r.name || r.characterId }}</strong>
                            <span class="votes">{{ r.votes }}</span>
                        </div>
                        <div class="bar" :style="{ width: Math.min(300, (r.votes || 0) * 10) + 'px' }"></div>
                    </li>
                </ul>
            </div>
            <div v-else class="empty-report">No votes yet.</div>
        </section>
    </div>
</template>

<style scoped>
.chart-container {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
}
.character-item { display:flex; gap:12px; margin:12px 0; align-items:flex-start }
.avatar { width:64px; height:64px; object-fit:cover; border-radius:8px }
.avatar-sm { width:32px; height:32px; object-fit:cover; border-radius:6px; margin-right:8px }
.info { flex:1 }
.vote-btn { margin-top:8px; padding:6px 10px; background:#007bff; color:#fff; border:none; border-radius:4px; cursor:pointer }
.vote-btn:hover { opacity:0.9 }
.report-section { margin-top:20px }
.report-item { display:flex; align-items:center; gap:12px; margin:8px 0 }
.r-info { display:flex; align-items:center; gap:8px; width:220px }
.bar { height:14px; background:#4caf50; border-radius:6px }
.votes { margin-left:8px; color:#333 }
</style>
<script setup lang="ts">
import { onMounted, ref } from 'vue';

const isBarSelected = ref(true);
const isTableSelected = ref(false);


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
        <canvas id="chart"></canvas>
    </div>
</template>

<style scoped>
.chart-container {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
}
</style>
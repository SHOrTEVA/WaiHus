<template>
    <nav class="navbar">
        <ul class="navbar-menu">
            <li class="navbar-item" v-for="item in menuItems" :key="item.text">
                <a :href="item.link" class="navbar-link">{{ item.text }}</a>
            </li>
        </ul>
    </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const menuItems = ref([
    { text: "Home", link: "/" },
    { text: "About", link: "/about" },
    { text: "Contact", link: "/contact" },
]);

</script>

<style scoped>
.navbar {
    background-color: #333;
    padding: 1rem;
}

.navbar-menu {
    list-style: none;
    display: flex;
    margin: 0;
    padding: 0;
}

.navbar-item {
    margin-right: 1rem;
}

.navbar-link {
    color: white;
    text-decoration: none;
}

.navbar-link:hover {
    text-decoration: underline;
}
</style>
<template>
  <v-navigation-drawer
    permanent
    expand-on-hover
    rail
    :rail-width="64"
    :width="200"
    rounded="xl"
    class="admin-sidebar custom-border dark-theme"
    color="grey-darken-3"
  >
      <v-list>
        <v-list-item class="my-2">
          <template #prepend>
            <v-avatar size="32">
              <v-icon color="blue" class="ml-1">mdi-all-inclusive</v-icon>
            </v-avatar>
          </template>
          <v-list-item-title class="ml-3 py-2 font-weight-bold text-h6">LingXi</v-list-item-title>
        </v-list-item>
        <v-divider class="my-1"></v-divider>

        <v-list density="compact" nav>
          <v-list-item to="/admin/user" prepend-icon="mdi-account-outline" class="ml-1" title="用户管理"></v-list-item>
          <v-list-item to="/admin/message" prepend-icon="mdi-message-text-outline" class="ml-1" title="公告"></v-list-item>
        </v-list>
      </v-list>
      <template #append>
        <v-divider class="my-1"></v-divider>
        <v-list density="compact" nav class="my-2">
          <v-list-item
            to="/"
            prepend-icon="mdi-logout"
            class="ml-1"
            title="用户视图"
          ></v-list-item>
        </v-list>
      </template>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>

</script>

<style scoped>
.admin-sidebar.dark-theme {
  --v-theme-background: #1e1e1e;
  --v-theme-surface: #2d2d2d;
}

.custom-border {
  border-top-right-radius: 16px !important;
  border-bottom-right-radius: 16px !important;
  border-top-left-radius: 0 !important;
  border-bottom-left-radius: 0 !important;
}

.v-list-item--active {
  background: rgba(66, 165, 245, 0.15) !important;
  border-left: 3px solid #42a5f5;
}

.v-list-item--active .v-icon {
  color: #42a5f5 !important;
}

.v-list-item--active .v-list-item-title {
  color: #42a5f5 !important;
  font-weight: 500;
}

.v-list-item:hover {
  background: rgba(255, 255, 255, 0.05) !important;
}

.admin-sidebar {
  height: 100vh;
  position: sticky;
  top: 0;
  transition: all 0.3s ease;
}

.v-icon {
  color: #58aef0 !important;
}
</style>

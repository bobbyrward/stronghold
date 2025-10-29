<script setup lang="ts">
defineProps<{
  id: string
  label?: string
  modelValue: number | string
  options: { value: number | string; label: string }[]
  placeholder?: string
  required?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number | string]
}>()

function onChange(event: Event) {
  const target = event.target as HTMLSelectElement
  const value = target.value
  // Convert to number if the original options use numbers
  emit('update:modelValue', isNaN(Number(value)) ? value : Number(value))
}
</script>

<template>
  <div class="mb-3">
    <label v-if="label" :for="id" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>
    <select
      :id="id"
      :value="modelValue"
      :required="required"
      class="form-select"
      @change="onChange"
    >
      <option v-if="placeholder" value="" disabled>
        {{ placeholder }}
      </option>
      <option
        v-for="option in options"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </option>
    </select>
  </div>
</template>

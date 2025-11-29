import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import DataTable from '../DataTable.vue'

describe('DataTable', () => {
    const columns = [
        { key: 'id', label: 'ID', editable: false },
        { key: 'name', label: 'Name', editable: true, type: 'text' as const }
    ]

    const data = [
        { id: 1, name: 'Item 1' },
        { id: 2, name: 'Item 2' }
    ]

    const defaultProps = {
        columns,
        data,
        loading: false,
        editable: true,
        onSave: vi.fn().mockResolvedValue(undefined),
        onDelete: vi.fn().mockResolvedValue(undefined)
    }

    it('renders columns and rows correctly', () => {
        const wrapper = mount(DataTable, {
            props: defaultProps
        })

        const headers = wrapper.findAll('th')
        expect(headers.length).toBe(3) // 2 columns + actions

        const rows = wrapper.findAll('tbody tr')
        expect(rows.length).toBe(2)

        expect(wrapper.text()).toContain('Item 1')
        expect(wrapper.text()).toContain('Item 2')
    })

    it('shows loading spinner when loading=true', () => {
        const wrapper = mount(DataTable, {
            props: {
                ...defaultProps,
                data: [],
                loading: true
            }
        })

        expect(wrapper.find('.spinner-border').exists()).toBe(true)
    })

    it('shows empty state message when data is empty', () => {
        const wrapper = mount(DataTable, {
            props: {
                ...defaultProps,
                data: [],
                loading: false
            }
        })

        expect(wrapper.text()).toContain('No data available')
    })

    it('clicking edit button enters edit mode', async () => {
        const wrapper = mount(DataTable, {
            props: defaultProps
        })

        const editButton = wrapper.find('[title="Edit"]')
        await editButton.trigger('click')

        expect(wrapper.find('input').exists()).toBe(true)
    })

    it('save button calls onSave with correct data', async () => {
        const onSave = vi.fn().mockResolvedValue(undefined)
        const wrapper = mount(DataTable, {
            props: {
                ...defaultProps,
                onSave
            }
        })

        // Enter edit mode
        const editButton = wrapper.find('[title="Edit"]')
        await editButton.trigger('click')

        // Click save
        const saveButton = wrapper.find('[title="Save"]')
        await saveButton.trigger('click')

        expect(onSave).toHaveBeenCalledWith(
            expect.objectContaining({ id: 1, name: 'Item 1' }),
            false
        )
    })

    it('cancel button exits edit mode', async () => {
        const wrapper = mount(DataTable, {
            props: defaultProps
        })

        // Enter edit mode
        const editButton = wrapper.find('[title="Edit"]')
        await editButton.trigger('click')

        expect(wrapper.find('input').exists()).toBe(true)

        // Click cancel
        const cancelButton = wrapper.find('[title="Cancel"]')
        await cancelButton.trigger('click')

        expect(wrapper.find('input').exists()).toBe(false)
    })

    it('delete button calls onDelete with item id', async () => {
        const onDelete = vi.fn().mockResolvedValue(undefined)

        // Mock window.confirm
        globalThis.confirm = vi.fn().mockReturnValue(true)

        const wrapper = mount(DataTable, {
            props: {
                ...defaultProps,
                onDelete
            }
        })

        const deleteButton = wrapper.find('[title="Delete"]')
        await deleteButton.trigger('click')

        expect(onDelete).toHaveBeenCalledWith(1)
    })
})

import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ConfirmDialog from '../ConfirmDialog.vue'

describe('ConfirmDialog', () => {
  const defaultProps = {
    show: true,
    title: 'Confirm Action',
    message: 'Are you sure?'
  }

  it('renders when show=true', () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    expect(wrapper.find('.modal').exists()).toBe(true)
    expect(wrapper.text()).toContain('Confirm Action')
    expect(wrapper.text()).toContain('Are you sure?')
  })

  it('hidden when show=false', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        ...defaultProps,
        show: false
      },
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    expect(wrapper.find('.modal').exists()).toBe(false)
  })

  it('emits confirm when confirm button clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    const confirmButton = wrapper.find('.btn-danger')
    await confirmButton.trigger('click')

    expect(wrapper.emitted('confirm')).toBeTruthy()
    expect(wrapper.emitted('confirm')?.length).toBe(1)
  })

  it('emits cancel when cancel button clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    const cancelButton = wrapper.find('.btn-secondary')
    await cancelButton.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
    expect(wrapper.emitted('cancel')?.length).toBe(1)
  })

  it('emits cancel when backdrop clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    const modal = wrapper.find('.modal')
    await modal.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('uses custom button text when provided', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        ...defaultProps,
        confirmText: 'Delete',
        cancelText: 'Keep'
      },
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    expect(wrapper.text()).toContain('Delete')
    expect(wrapper.text()).toContain('Keep')
  })

  it('uses specified variant for confirm button', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        ...defaultProps,
        variant: 'warning'
      },
      global: {
        stubs: {
          Teleport: true
        }
      }
    })

    expect(wrapper.find('.btn-warning').exists()).toBe(true)
  })
})

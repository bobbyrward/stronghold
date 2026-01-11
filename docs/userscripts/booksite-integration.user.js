// ==UserScript==
// @name         Stronghold book site Integration
// @namespace    https://github.com/bobbyrward/stronghold
// @version      1.0.0
// @description  Add torrents from BookSite to Stronghold for automatic book import
// @match        https://www.*.net/t/*
// @grant        GM_xmlhttpRequest
// @connect      stronghold.home.ohnozombi.es
// @connect      localhost
// @run-at       document-idle
// @updateURL    https://raw.githubusercontent.com/bobbyrward/stronghold/main/docs/userscripts/booksite-integration.user.js
// @downloadURL  https://raw.githubusercontent.com/bobbyrward/stronghold/main/docs/userscripts/booksite-integration.user.js
// ==/UserScript==

(function () {
  'use strict';

  // ===== CONFIGURATION =====
  const CONFIG = {
    apiBaseUrl: 'https://stronghold.home.ohnozombi.es/api',
    buttonSelector: '#download .torDetInnerBottom',
    timeout: 10000,
    debug: false
  };

  // ===== CATEGORIES =====
  const CATEGORIES = [
    'books',
    'audiobooks',
    'personal-books',
    'personal-audiobooks',
    'general-books',
    'general-audiobooks',
    'kids-books',
    'kids-audiobooks'
  ];

  // ===== STATE =====
  let isModalOpen = false;
  let isSubmitting = false;
  let modalElement = null;

  // ===== UTILITY FUNCTIONS =====

  /**
   * Make an API request using GM_xmlhttpRequest
   * @param {Object} options - Request options
   * @returns {Promise} Promise resolving to response data
   */
  function makeApiRequest(options) {
    return new Promise((resolve, reject) => {
      GM_xmlhttpRequest({
        method: options.method || 'GET',
        url: options.url,
        headers: {
          'Content-Type': 'application/json',
          ...options.headers
        },
        data: options.data ? JSON.stringify(options.data) : undefined,
        timeout: CONFIG.timeout,

        onload: function (response) {
          if (response.status >= 200 && response.status < 300) {
            try {
              const data = response.responseText ?
                JSON.parse(response.responseText) : null;
              resolve({ success: true, data, status: response.status });
            } catch (e) {
              reject(new Error('Invalid JSON response'));
            }
          } else {
            let errorMessage = `HTTP ${response.status}`;
            try {
              const errorData = JSON.parse(response.responseText);
              errorMessage = errorData.error || errorMessage;
            } catch (e) {
              errorMessage = response.statusText || errorMessage;
            }
            reject(new Error(errorMessage));
          }
        },

        onerror: function () {
          reject(new Error('Network error - cannot reach Stronghold API'));
        },

        ontimeout: function () {
          reject(new Error('Request timeout - Stronghold API not responding'));
        }
      });
    });
  }

  /**
   * Sanitize input to prevent XSS
   * @param {string} input - Input string
   * @returns {string} Sanitized string
   */
  function sanitizeInput(input) {
    if (typeof input !== 'string') {
      return '';
    }
    const div = document.createElement('div');
    div.textContent = input;
    return div.innerHTML;
  }

  /**
   * Log debug message
   * @param {string} message - Debug message
   * @param  {...any} args - Additional arguments
   */
  function debug(message, ...args) {
    if (CONFIG.debug) {
      console.log(`[Stronghold] ${message}`, ...args);
    }
  }

  // ===== UI COMPONENT FUNCTIONS =====

  /**
   * Create the Stronghold button
   * @returns {HTMLElement} Button element
   */
  function createButton() {
    const button = document.createElement('button');
    button.id = 'stronghold-dl-btn';
    button.textContent = '⚡ Download to Stronghold';
    button.title = 'Send this torrent to Stronghold for automatic import';

    // Inline styles (required due to book sites' CSP)
    Object.assign(button.style, {
      backgroundColor: '#4CAF50',
      color: 'white',
      padding: '10px 20px',
      margin: '8px',
      border: 'none',
      borderRadius: '4px',
      cursor: 'pointer',
      fontSize: '14px',
      fontWeight: 'bold',
      marginTop: '10px',
      transition: 'background-color 0.3s ease',
      boxShadow: '0 2px 4px rgba(0,0,0,0.2)'
    });

    // Hover effects
    button.addEventListener('mouseenter', () => {
      button.style.backgroundColor = '#45a049';
    });
    button.addEventListener('mouseleave', () => {
      button.style.backgroundColor = '#4CAF50';
    });

    return button;
  }

  /**
   * Inject the button into the page
   * @returns {boolean} True if successful, false otherwise
   */
  function injectButton() {
    const targetElement = document.querySelector(CONFIG.buttonSelector);

    if (!targetElement) {
      console.error('[Stronghold] Target element not found:', CONFIG.buttonSelector);
      return false;
    }

    // Prevent duplicate injection
    if (document.getElementById('stronghold-dl-btn')) {
      debug('Button already exists, skipping injection');
      return true;
    }

    const button = createButton();
    button.addEventListener('click', handleButtonClick);

    // Insert after target element
    targetElement.parentNode.insertBefore(button, targetElement.nextSibling);

    debug('Button injected successfully');
    return true;
  }

  /**
   * Create the modal dialog
   * @returns {HTMLElement} Modal element
   */
  function createModal() {
    const modal = document.createElement('div');
    modal.id = 'stronghold-modal';

    modal.innerHTML = `
            <div class="stronghold-modal-overlay"></div>
            <div class="stronghold-modal-content">
                <div class="stronghold-modal-header">
                    <h2>Download to Stronghold</h2>
                    <button class="stronghold-modal-close" aria-label="Close">&times;</button>
                </div>
                <div class="stronghold-modal-body">
                    <div class="stronghold-form-group">
                        <label for="stronghold-category-select">Select Category:</label>
                        <select id="stronghold-category-select" class="stronghold-select">
                            <option value="" disabled selected>-- Select a category --</option>
                        </select>
                    </div>
                    <div id="stronghold-torrent-info" class="stronghold-info-box">
                        <p><strong>Torrent:</strong> <span id="torrent-name-display">Loading...</span></p>
                        <p><strong>ID:</strong> <span id="torrent-id-display">Loading...</span></p>
                    </div>
                </div>
                <div class="stronghold-modal-footer">
                    <button id="stronghold-submit-btn" class="stronghold-btn stronghold-btn-primary">
                        Submit
                    </button>
                    <button id="stronghold-cancel-btn" class="stronghold-btn stronghold-btn-secondary">
                        Cancel
                    </button>
                </div>
            </div>
        `;

    applyModalStyles(modal);
    setupModalEventListeners(modal);

    return modal;
  }

  /**
   * Apply inline styles to modal elements
   * @param {HTMLElement} modal - Modal element
   */
  function applyModalStyles(modal) {
    // Overlay (backdrop)
    const overlay = modal.querySelector('.stronghold-modal-overlay');
    Object.assign(overlay.style, {
      position: 'fixed',
      top: '0',
      left: '0',
      width: '100%',
      height: '100%',
      backgroundColor: 'rgba(0, 0, 0, 0.6)',
      zIndex: '9998',
      display: 'block'
    });

    // Modal content box
    const content = modal.querySelector('.stronghold-modal-content');
    Object.assign(content.style, {
      position: 'fixed',
      top: '50%',
      left: '50%',
      transform: 'translate(-50%, -50%)',
      backgroundColor: '#202020',
      borderRadius: '8px',
      boxShadow: '0 4px 20px rgba(0,0,0,0.3)',
      zIndex: '9999',
      minWidth: '400px',
      maxWidth: '600px',
      padding: '0'
    });

    // Header
    const header = modal.querySelector('.stronghold-modal-header');
    Object.assign(header.style, {
      padding: '20px 24px',
      borderBottom: '1px solid #e0e0e0',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center'
    });

    const h2 = header.querySelector('h2');
    Object.assign(h2.style, {
      margin: '0',
      fontSize: '20px',
      fontWeight: 'bold',
      color: 'rgb(188,188,188)'
    });

    const closeBtn = header.querySelector('.stronghold-modal-close');
    Object.assign(closeBtn.style, {
      background: 'none',
      border: 'none',
      fontSize: '28px',
      cursor: 'pointer',
      color: '#999',
      padding: '0',
      width: '30px',
      height: '30px',
      lineHeight: '1'
    });

    // Body
    const body = modal.querySelector('.stronghold-modal-body');
    Object.assign(body.style, {
      padding: '24px'
    });

    // Form group
    const formGroup = modal.querySelector('.stronghold-form-group');
    Object.assign(formGroup.style, {
      marginBottom: '20px'
    });

    const label = formGroup.querySelector('label');
    Object.assign(label.style, {
      display: 'block',
      marginBottom: '8px',
      fontWeight: 'bold',
      color: 'rgb(188,188,188)',
      fontSize: '14px'
    });

    const select = formGroup.querySelector('select');
    Object.assign(select.style, {
      width: '100%',
      padding: '10px',
      fontSize: '14px',
      border: '1px solid #ddd',
      borderRadius: '4px',
      backgroundColor: '#fff',
      cursor: 'pointer'
    });

    // Info box
    const infoBox = modal.querySelector('.stronghold-info-box');
    Object.assign(infoBox.style, {
      backgroundColor: '#f5f5f5',
      padding: '16px',
      borderRadius: '4px',
      fontSize: '14px'
    });

    const infoPs = infoBox.querySelectorAll('p');
    infoPs.forEach(p => {
      Object.assign(p.style, {
        margin: '8px 0',
        color: '#555'
      });
    });

    // Footer
    const footer = modal.querySelector('.stronghold-modal-footer');
    Object.assign(footer.style, {
      padding: '16px 24px',
      borderTop: '1px solid #e0e0e0',
      display: 'flex',
      justifyContent: 'flex-end',
      gap: '12px'
    });

    // Buttons
    const buttons = footer.querySelectorAll('.stronghold-btn');
    buttons.forEach(btn => {
      Object.assign(btn.style, {
        padding: '10px 20px',
        fontSize: '14px',
        fontWeight: 'bold',
        border: 'none',
        borderRadius: '4px',
        cursor: 'pointer',
        transition: 'background-color 0.3s ease'
      });
    });

    const submitBtn = footer.querySelector('.stronghold-btn-primary');
    Object.assign(submitBtn.style, {
      backgroundColor: '#4CAF50',
      color: 'white'
    });
    submitBtn.addEventListener('mouseenter', () => {
      submitBtn.style.backgroundColor = '#45a049';
    });
    submitBtn.addEventListener('mouseleave', () => {
      submitBtn.style.backgroundColor = '#4CAF50';
    });

    const cancelBtn = footer.querySelector('.stronghold-btn-secondary');
    Object.assign(cancelBtn.style, {
      backgroundColor: '#f0f0f0',
      color: '#333'
    });
    cancelBtn.addEventListener('mouseenter', () => {
      cancelBtn.style.backgroundColor = '#e0e0e0';
    });
    cancelBtn.addEventListener('mouseleave', () => {
      cancelBtn.style.backgroundColor = '#f0f0f0';
    });
  }

  /**
   * Setup event listeners for modal
   * @param {HTMLElement} modal - Modal element
   */
  function setupModalEventListeners(modal) {
    // Close button
    const closeBtn = modal.querySelector('.stronghold-modal-close');
    closeBtn.addEventListener('click', handleModalCancel);

    // Cancel button
    const cancelBtn = modal.querySelector('#stronghold-cancel-btn');
    cancelBtn.addEventListener('click', handleModalCancel);

    // Submit button
    const submitBtn = modal.querySelector('#stronghold-submit-btn');
    submitBtn.addEventListener('click', handleModalSubmit);

    // Overlay click
    const overlay = modal.querySelector('.stronghold-modal-overlay');
    overlay.addEventListener('click', handleModalCancel);

    // ESC key
    document.addEventListener('keydown', handleEscKey);
  }

  /**
   * Show the modal dialog
   */
  function showModal() {
    if (!modalElement) {
      modalElement = createModal();
      document.body.appendChild(modalElement);
    }

    // Populate category select
    populateCategorySelect();

    // Update torrent info
    updateTorrentInfo();

    // Show modal
    modalElement.style.display = 'block';
    isModalOpen = true;

    debug('Modal opened');
  }

  /**
   * Hide the modal dialog
   */
  function hideModal() {
    if (modalElement) {
      modalElement.remove();
      modalElement = null;
    }

    // Remove ESC key listener to prevent memory leak
    document.removeEventListener('keydown', handleEscKey);

    isModalOpen = false;
    debug('Modal closed');
  }

  /**
   * Show feedback toast notification
   * @param {boolean} success - True for success, false for error
   * @param {string} message - Message to display
   */
  function showFeedback(success, message) {
    // Remove existing toast
    const existingToast = document.getElementById('stronghold-toast');
    if (existingToast) {
      existingToast.remove();
    }

    // Create toast
    const toast = document.createElement('div');
    toast.id = 'stronghold-toast';
    const icon = success ? '✓' : '✗';

    toast.innerHTML = `
            <span class="toast-icon">${icon}</span>
            <span class="toast-message">${sanitizeInput(message)}</span>
        `;

    // Styling
    Object.assign(toast.style, {
      position: 'fixed',
      top: '20px',
      right: '20px',
      backgroundColor: success ? '#4CAF50' : '#f44336',
      color: 'white',
      padding: '16px 24px',
      borderRadius: '4px',
      boxShadow: '0 4px 12px rgba(0,0,0,0.3)',
      zIndex: '10000',
      display: 'flex',
      alignItems: 'center',
      gap: '12px',
      fontSize: '14px',
      fontWeight: 'bold',
      maxWidth: '400px'
    });

    document.body.appendChild(toast);

    // Auto-remove after 5 seconds
    setTimeout(() => {
      if (toast.parentNode) {
        toast.remove();
      }
    }, 5000);

    debug(`Toast shown: ${success ? 'success' : 'error'} - ${message}`);
  }

  // ===== DATA FUNCTIONS =====

  /**
   * Extract torrent data from the page
   * @returns {Object} Torrent data
   */
  function extractTorrentData() {
    const torrentData = {
      torrentId: window.location.pathname.match(/\/t\/(\d+)/)[1],
      torrentName: document.querySelector(".TorrentTitle").textContent.trim(),
      torrentUrl: document.querySelector("#tddl").href,
      category: null
    };

    debug('Extracted torrent data:', torrentData);
    return torrentData;
  }

  /**
   * Validate torrent data
   * @param {Object} torrentData - Torrent data to validate
   * @returns {Object} Validation result {valid: boolean, errors: string[]}
   */
  function validateTorrentData(torrentData) {
    const errors = [];

    if (!torrentData.torrentId) {
      errors.push('Torrent ID could not be determined');
    }
    if (!torrentData.torrentName) {
      errors.push('Torrent name could not be found');
    }
    if (!torrentData.torrentUrl) {
      errors.push('Torrent url could not be found');
    }

    if (!torrentData.category) {
      errors.push('Please select a category');
    }

    return {
      valid: errors.length === 0,
      errors: errors
    };
  }

  /**
   * Populate the category select dropdown
   */
  function populateCategorySelect() {
    const select = document.getElementById('stronghold-category-select');
    if (!select) return;

    // Clear existing options except the first one
    while (select.options.length > 1) {
      select.remove(1);
    }

    // Add categories
    CATEGORIES.forEach(category => {
      const option = document.createElement('option');
      option.value = category;
      option.textContent = category;
      select.appendChild(option);
    });

    debug('Category select populated');
  }

  /**
   * Update torrent info display in modal
   */
  function updateTorrentInfo() {
    const torrentData = extractTorrentData();

    const nameDisplay = document.getElementById('torrent-name-display');
    const idDisplay = document.getElementById('torrent-id-display');

    if (nameDisplay) {
      nameDisplay.textContent = torrentData.torrentName || 'Unknown';
    }
    if (idDisplay) {
      idDisplay.textContent = torrentData.torrentId || 'Unknown';
    }

    debug('Torrent info updated');
  }

  // ===== API FUNCTIONS =====

  /**
   * Submit torrent download request to Stronghold API
   * @param {Object} torrentData - Torrent data to submit
   * @returns {Promise} Promise resolving to response data
   */
  async function submitTorrentDownload(torrentData) {
    try {
      const result = await makeApiRequest({
        method: 'POST',
        url: `${CONFIG.apiBaseUrl}/book-torrent-dl`,
        data: {
          category: torrentData.category,
          torrent_id: torrentData.torrentId,
          torrent_name: torrentData.torrentName,
          torrent_url: torrentData.torrentUrl
        }
      });

      if (result.success) {
        debug('Torrent submitted successfully:', result.data);
        return result.data;
      } else {
        throw new Error('Submission failed');
      }
    } catch (error) {
      console.error('[Stronghold] Submission error:', error);
      throw error;
    }
  }

  // ===== EVENT HANDLERS =====

  /**
   * Handle button click event
   * @param {Event} event - Click event
   */
  function handleButtonClick(event) {
    event.preventDefault();
    event.stopPropagation();

    // Prevent multiple modal instances
    if (isModalOpen) {
      debug('Modal already open');
      return;
    }

    // Prevent clicks during submission
    if (isSubmitting) {
      debug('Submission in progress');
      return;
    }

    showModal();
  }

  /**
   * Handle modal submit event
   * @param {Event} event - Click event
   */
  async function handleModalSubmit(event) {
    event.preventDefault();

    // Prevent double submission
    if (isSubmitting) {
      return;
    }

    // Extract torrent data
    const torrentData = extractTorrentData();

    // Get selected category
    const select = document.getElementById('stronghold-category-select');
    torrentData.category = select ? select.value : null;

    // Validate
    const validation = validateTorrentData(torrentData);
    if (!validation.valid) {
      showFeedback(false, validation.errors.join('. '));
      return;
    }

    // Update submit button
    const submitButton = document.getElementById('stronghold-submit-btn');
    const originalText = submitButton.textContent;
    submitButton.textContent = 'Submitting...';
    submitButton.disabled = true;
    isSubmitting = true;

    try {
      await submitTorrentDownload(torrentData);
      hideModal();
      showFeedback(true, 'Torrent successfully added to Stronghold!');
    } catch (error) {
      let errorMessage = 'Failed to submit torrent';

      if (error.message.includes('Network error')) {
        errorMessage = 'Cannot connect to Stronghold API. Is the server running?';
      } else if (error.message.includes('timeout')) {
        errorMessage = 'Request timed out. Please try again.';
      } else if (error.message) {
        errorMessage = error.message;
      }

      showFeedback(false, errorMessage);

      // Re-enable submit button on error
      submitButton.textContent = originalText;
      submitButton.disabled = false;
    } finally {
      isSubmitting = false;
    }
  }

  /**
   * Handle modal cancel event
   * @param {Event} event - Click event
   */
  function handleModalCancel(event) {
    event.preventDefault();
    hideModal();
  }

  /**
   * Handle ESC key press
   * @param {KeyboardEvent} event - Keyboard event
   */
  function handleEscKey(event) {
    if (event.key === 'Escape' && isModalOpen) {
      hideModal();
    }
  }

  // ===== INITIALIZATION =====

  /**
   * Initialize the userscript
   */
  function init() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', initializeScript);
    } else {
      initializeScript();
    }
  }

  /**
   * Initialize script after DOM is ready
   */
  function initializeScript() {
    console.log('[Stronghold] Booksite Integration userscript loaded');
    injectButton();
  }

  // Start
  init();
})();

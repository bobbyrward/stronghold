// ==UserScript==
// @name         Stronghold Browse Integration
// @namespace    https://github.com/bobbyrward/stronghold
// @version      1.0.0
// @description  Add torrents from browse page to Stronghold for automatic book import
// @match        https://www.*.net/tor/browse.php*
// @grant        GM_xmlhttpRequest
// @connect      stronghold.home.ohnozombi.es
// @connect      localhost
// @run-at       document-idle
// @updateURL    https://raw.githubusercontent.com/bobbyrward/stronghold/main/docs/userscripts/browse-integration.user.js
// @downloadURL  https://raw.githubusercontent.com/bobbyrward/stronghold/main/docs/userscripts/browse-integration.user.js
// ==/UserScript==

(function () {
  'use strict';

  // ===== CONFIGURATION =====
  const CONFIG = {
    apiBaseUrl: 'https://stronghold.home.ohnozombi.es/api',
    tableSelector: 'table.newTorTable',
    timeout: 10000,
    debug: false
  };

  // ===== SVG ICON (data URL) =====
  const STRONGHOLD_ICON = 'data:image/svg+xml;base64,PHN2ZyBmaWxsPSIjNENBRjUwIiB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgd2lkdGg9IjI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgNDY2LjcxNSA0NjYuNzE1Ij48Zz48cGF0aCBkPSJNMzY4LjA4NiwyMDguNzQzVjg1Ljg0aDIwLjg1NVYwaC0zOC43MjN2NDMuODU3aC0yOS4xOFYwLjAwMWgtMzkuMDAydjQzLjg1NmgtMjkuMTc4VjAuMDAxaC0zOS4wMDR2NDMuODU2aC0yOS4xNzhWMC4wMDFoLTM5LjAwM3Y0My44NTZoLTI5LjE3OFYwLjAwMUg3Ny43NzN2ODUuODRoMjAuODU1djEyMi45MDJINjUuNXY5Mi4xMTFoMjIuNTU3djE2NS44NjFoNTEuMjk1di02NC4xODloNjQuODU0djY0LjE4OWg1OC4zMDJ2LTY0LjE4OWg2NC44NTR2NjQuMTg5aDUxLjI5NVYzMDAuODU0aDIyLjU1OXYtOTIuMTExSDM2OC4wODZ6IE0yMDYuMDEzLDE2OS45OTl2MzguNzQ0aC01OS4xNzh2LTM4Ljc0NEgyMDYuMDEzeiBNMzE5Ljg3OSwxNjkuOTk5djM4Ljc0NGgtNTkuMTh2LTM4Ljc0NEgzMTkuODc5eiIvPjwvZz48L3N2Zz4=';

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
  let currentTorrentData = null;

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
   * Create a Stronghold button for a torrent row
   * @param {Object} torrentData - Data for the torrent
   * @returns {HTMLElement} Button element
   */
  function createButton(torrentData) {
    const button = document.createElement('a');
    button.className = 'stronghold-dl-btn';
    button.title = 'Send to Stronghold';
    button.setAttribute('role', 'button');

    const img = document.createElement('img');
    img.src = STRONGHOLD_ICON;
    img.alt = 'Stronghold';
    Object.assign(img.style, {
      width: '16px',
      height: '16px',
      verticalAlign: 'middle'
    });

    button.appendChild(img);

    // Store torrent data on the button
    button.dataset.torrentId = torrentData.torrentId;
    button.dataset.torrentName = torrentData.torrentName;
    button.dataset.torrentUrl = torrentData.torrentUrl;

    // Inline styles
    Object.assign(button.style, {
      cursor: 'pointer',
      marginLeft: '8px',
      display: 'inline-block',
      padding: '2px',
      borderRadius: '3px',
      transition: 'background-color 0.2s ease'
    });

    // Hover effects
    button.addEventListener('mouseenter', () => {
      button.style.backgroundColor = 'rgba(76, 175, 80, 0.2)';
    });
    button.addEventListener('mouseleave', () => {
      button.style.backgroundColor = 'transparent';
    });

    return button;
  }

  /**
   * Inject buttons into all torrent rows
   */
  function injectButtons() {
    const table = document.querySelector(CONFIG.tableSelector);

    if (!table) {
      debug('Torrent table not found');
      return false;
    }

    const rows = table.querySelectorAll('tr[id^="tdr-"]');

    if (rows.length === 0) {
      debug('No torrent rows found');
      return false;
    }

    rows.forEach(row => {
      // Skip if button already exists
      if (row.querySelector('.stronghold-dl-btn')) {
        return;
      }

      const torTitle = row.querySelector('a.torTitle');
      const directDownload = row.querySelector('a.directDownload');

      if (!torTitle || !directDownload) {
        debug('Missing torTitle or directDownload in row', row.id);
        return;
      }

      // Extract torrent data
      const torrentData = {
        torrentId: torTitle.href.match(/\/t\/(\d+)/)?.[1],
        torrentName: torTitle.textContent.trim(),
        torrentUrl: directDownload.href
      };

      if (!torrentData.torrentId) {
        debug('Could not extract torrent ID from', torTitle.href);
        return;
      }

      // Create and inject button
      const button = createButton(torrentData);
      button.addEventListener('click', handleButtonClick);

      // Insert after the last link in the download cell
      const downloadCell = directDownload.parentElement;
      downloadCell.appendChild(button);
    });

    debug(`Injected buttons into ${rows.length} rows`);
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
   * @param {Object} torrentData - Torrent data to display
   */
  function showModal(torrentData) {
    if (!modalElement) {
      modalElement = createModal();
      document.body.appendChild(modalElement);
    }

    currentTorrentData = torrentData;

    // Populate category select
    populateCategorySelect();

    // Update torrent info
    updateTorrentInfo(torrentData);

    // Show modal
    modalElement.style.display = 'block';
    isModalOpen = true;

    debug('Modal opened for torrent:', torrentData.torrentName);
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
    currentTorrentData = null;
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
   * @param {Object} torrentData - Torrent data to display
   */
  function updateTorrentInfo(torrentData) {
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

    // Get torrent data from button's data attributes
    const button = event.currentTarget;
    const torrentData = {
      torrentId: button.dataset.torrentId,
      torrentName: button.dataset.torrentName,
      torrentUrl: button.dataset.torrentUrl
    };

    showModal(torrentData);
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

    if (!currentTorrentData) {
      showFeedback(false, 'No torrent data available');
      return;
    }

    // Get selected category
    const select = document.getElementById('stronghold-category-select');
    currentTorrentData.category = select ? select.value : null;

    // Validate
    const validation = validateTorrentData(currentTorrentData);
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
      await submitTorrentDownload(currentTorrentData);
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
   * Wait for an element to appear in the DOM
   * @param {string} selector - CSS selector
   * @param {number} timeout - Timeout in milliseconds
   * @returns {Promise<Element>} Promise resolving to the element
   */
  function waitForElement(selector, timeout = 10000) {
    return new Promise((resolve, reject) => {
      const element = document.querySelector(selector);
      if (element) {
        resolve(element);
        return;
      }

      const observer = new MutationObserver((mutations, obs) => {
        const element = document.querySelector(selector);
        if (element) {
          obs.disconnect();
          resolve(element);
        }
      });

      observer.observe(document.body, { childList: true, subtree: true });

      setTimeout(() => {
        observer.disconnect();
        reject(new Error(`Element ${selector} not found within ${timeout}ms`));
      }, timeout);
    });
  }

  /**
   * Initialize script after DOM is ready
   */
  async function initializeScript() {
    console.log('[Stronghold] Browse Integration userscript loaded');

    try {
      // Wait for the table to appear (it may be loaded dynamically)
      const table = await waitForElement(CONFIG.tableSelector);
      debug('Table found, injecting buttons');
      injectButtons();

      // Re-inject buttons when page content changes (e.g., AJAX pagination)
      const observer = new MutationObserver((mutations) => {
        for (const mutation of mutations) {
          if (mutation.addedNodes.length > 0) {
            // Check if new torrent rows were added
            const hasNewRows = Array.from(mutation.addedNodes).some(node =>
              node.nodeType === 1 && (
                node.matches?.('tr[id^="tdr-"]') ||
                node.querySelector?.('tr[id^="tdr-"]')
              )
            );
            if (hasNewRows) {
              debug('New rows detected, injecting buttons');
              injectButtons();
            }
          }
        }
      });

      observer.observe(table, { childList: true, subtree: true });
    } catch (error) {
      console.log('[Stronghold] No torrent table found on this page');
    }
  }

  // Start
  init();
})();

/**
 * Recipe UI Enhancements
 * Handles collapsible menus and scroll animations
 */

document.addEventListener('DOMContentLoaded', () => {
    initCollapsibles();
    initScrollAnimations();
});

/**
 * Initialize collapsible menus (accordions)
 */
function initCollapsibles() {
    const triggers = document.querySelectorAll('.faq-toggle, .collapsible-trigger');

    triggers.forEach(trigger => {
        trigger.addEventListener('click', () => {
            const isExpanded = trigger.getAttribute('aria-expanded') === 'true';
            const contentId = trigger.getAttribute('aria-controls') ||
                trigger.nextElementSibling?.id;

            if (!contentId) return;

            const content = document.getElementById(contentId) || trigger.nextElementSibling;

            // Toggle state
            trigger.setAttribute('aria-expanded', !isExpanded);

            if (content) {
                if (!isExpanded) {
                    content.style.display = 'block';
                    // Small delay to allow display:block to apply before opacity transition
                    requestAnimationFrame(() => {
                        content.classList.add('expanded');
                    });
                } else {
                    content.classList.remove('expanded');
                    // Wait for transition to finish before hiding
                    content.addEventListener('transitionend', function handler() {
                        if (!content.classList.contains('expanded')) {
                            content.style.display = 'none';
                        }
                        content.removeEventListener('transitionend', handler);
                    }, { once: true });
                }
            }

            // Rotate icon if present
            const icon = trigger.querySelector('.icon-arrow');
            if (icon) {
                icon.style.transform = isExpanded ? 'rotate(0deg)' : 'rotate(180deg)';
            }
        });
    });
}

/**
 * Initialize scroll animations using Intersection Observer
 */
function initScrollAnimations() {
    const animatedElements = document.querySelectorAll('[data-aos]');

    const observerOptions = {
        root: null,
        rootMargin: '0px',
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('aos-animate');
                // Optional: Stop observing once animated
                // observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    animatedElements.forEach(el => {
        observer.observe(el);
    });
}

// Smooth scroll behavior
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Add scroll animation for feature cards
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -100px 0px'
};

const observer = new IntersectionObserver(function(entries) {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Observe all cards
document.querySelectorAll('.feature-card, .arch-component, .use-case, .tech-item, .step').forEach(card => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(20px)';
    card.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    observer.observe(card);
});

// Mobile menu toggle (if needed)
const navLinks = document.querySelector('.nav-links');
let scrollPosition = 0;

window.addEventListener('scroll', () => {
    scrollPosition = window.scrollY;
});

// Add active class to nav links based on scroll position
function updateActiveNavLink() {
    const sections = document.querySelectorAll('section[id]');
    
    sections.forEach(section => {
        const sectionTop = section.offsetTop - 100;
        const sectionBottom = sectionTop + section.offsetHeight;
        
        if (scrollPosition >= sectionTop && scrollPosition < sectionBottom) {
            document.querySelectorAll('.nav-links a').forEach(link => {
                link.style.color = 'var(--text-dark)';
                if (link.getAttribute('href') === '#' + section.id) {
                    link.style.color = 'var(--primary)';
                }
            });
        }
    });
}

window.addEventListener('scroll', updateActiveNavLink);

// Update hero image on scroll
window.addEventListener('scroll', () => {
    const hero = document.querySelector('.hero-image');
    if (hero) {
        const scrolled = window.scrollY;
        hero.style.transform = `translateY(${scrolled * 0.5}px)`;
    }
});

console.log('Tech Debt Collector website loaded!');

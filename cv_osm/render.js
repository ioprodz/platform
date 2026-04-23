// CV renderer: reads ?lang & ?view from URL, fetches matching JSON,
// and hydrates the DOM with either the Brief or Extended view.

(function () {
  const SUPPORTED_LANGS = ['en', 'fr'];
  const SUPPORTED_VIEWS = ['brief', 'extended'];

  const $ = (sel, root) => (root || document).querySelector(sel);

  function esc(value) {
    if (value == null) return '';
    return String(value).replace(/[&<>"']/g, (c) => ({
      '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
    }[c]));
  }

  function parseParams() {
    const p = new URLSearchParams(location.search);
    let lang = p.get('lang');
    if (!SUPPORTED_LANGS.includes(lang)) lang = 'en';
    let view = p.get('view');
    if (!SUPPORTED_VIEWS.includes(view)) view = 'brief';
    return { lang, view };
  }

  function setUrl(lang, view) {
    const p = new URLSearchParams();
    p.set('lang', lang);
    p.set('view', view);
    history.pushState({ lang, view }, '', '?' + p.toString());
  }

  function track(event, params) {
    if (typeof window.plausible !== 'function') return;
    try { window.plausible(event, params); } catch (e) { /* ignore */ }
  }

  function domainOf(website) {
    if (!website) return null;
    try { return new URL(website).hostname.replace(/^www\./, ''); }
    catch (e) { return null; }
  }

  function logoImg(website, className) {
    const domain = domainOf(website);
    if (!domain) return '';
    const src = 'https://www.google.com/s2/favicons?domain=' + encodeURIComponent(domain) + '&sz=64';
    const img = '<img class="' + className + '" src="' + src +
      '" alt="" loading="lazy" referrerpolicy="no-referrer" ' +
      'onerror="this.parentElement&&this.parentElement.remove()">';
    return '<a class="' + className + '-link" href="' + esc(website) +
      '" target="_blank" rel="noopener noreferrer" aria-label="' + esc(domain) + '">' +
      img + '</a>';
  }

  function linkedText(text, website, className) {
    if (!website) return esc(text);
    return '<a class="' + className + '" href="' + esc(website) +
      '" target="_blank" rel="noopener noreferrer">' + esc(text) + '</a>';
  }

  const dataCache = {};
  async function loadData(lang) {
    if (dataCache[lang]) return dataCache[lang];
    const res = await fetch('data/' + lang + '.json', { cache: 'no-cache' });
    if (!res.ok) throw new Error('Failed to load data/' + lang + '.json');
    dataCache[lang] = await res.json();
    return dataCache[lang];
  }

  // ---------- Renderers ----------

  function renderHeader(meta) {
    return [
      '<div class="header">',
      '  <div class="header-left">',
      '    <h1>' + esc(meta.name) + '</h1>',
      '    <div class="title">' + esc(meta.title) + '</div>',
      '    <div class="experience">' + esc(meta.experienceSummary) + '</div>',
      '    <div class="contact-info">',
      '      <div>' + esc(meta.contact.phone) + '</div>',
      '      <div>' + esc(meta.contact.email) + '</div>',
      '      <div>' + esc(meta.contact.address) + '</div>',
      '    </div>',
      '  </div>',
      '  <div class="header-right">',
      '    <img src="' + esc(meta.avatar) + '" alt="' + esc(meta.name) + '">',
      '  </div>',
      '</div>'
    ].join('\n');
  }

  function renderBrief(data) {
    const meta = data.meta;
    const ui = data.ui;
    const params = parseParams();
    const extendedHref = './?lang=' + encodeURIComponent(params.lang) + '&view=extended';

    const skillsLines = data.skills.map(function (s) {
      return '<div class="brief-skills-line"><strong>' +
        esc(s.category) + '</strong> — ' + esc(s.items) + '</div>';
    }).join('');

    const briefJobs = data.experience.filter(function (e) { return e.showInBrief; }).map(function (e) {
      const companyInner = e.company
        ? (e.website
          ? '<a class="company company-link" href="' + esc(e.website) + '" target="_blank" rel="noopener noreferrer">' + esc(e.company) + '</a>'
          : '<span class="company">' + esc(e.company) + '</span>')
        : '';
      const companyPart = e.company ? ' · ' + companyInner : '';
      const periodPart = e.period
        ? '<span class="period">' + esc(e.period) + '</span>'
        : '';
      const line = e.briefLine
        ? '<div class="brief-job-line">' + esc(e.briefLine) + '</div>'
        : '';
      const logo = logoImg(e.website, 'brief-logo');
      return [
        '<div class="brief-job">',
        '  ' + logo,
        '  <div class="brief-job-body">',
        '    <div class="brief-job-head">' + esc(e.title) + companyPart + ' ' + periodPart + '</div>',
        '    ' + line,
        '  </div>',
        '</div>'
      ].join('\n');
    }).join('');

    const eduItems = data.education.map(function (e) {
      const schoolHtml = linkedText(e.school, e.website, 'school-link');
      return '<div class="brief-edu-item">' +
        '<span class="date">' + esc(e.date) + '</span>' +
        esc(e.degree) + ' — ' + schoolHtml +
        '</div>';
    }).join('');

    return [
      renderHeader(meta),
      '<p class="bio">' + esc(meta.bio) + '</p>',
      '<div class="brief-skills">' + skillsLines + '</div>',
      '<div class="section">',
      '  <h2 class="section-title">' + esc(ui.sectionExperience) + '</h2>',
      '  ' + briefJobs,
      '</div>',
      '<div class="section">',
      '  <h2 class="section-title">' + esc(ui.sectionEducation) + '</h2>',
      '  ' + eduItems,
      '</div>',
      '<div class="brief-footer">' +
        '<a class="brief-footer-link" href="' + esc(extendedHref) + '">' +
          esc(ui.briefFooterNote) +
        '</a>' +
      '</div>'
    ].join('\n');
  }

  function renderJob(job, ui) {
    const headBits = [];
    headBits.push('<div class="job-title">' + esc(job.title) + '</div>');

    if (job.company) {
      const locPart = job.location ? ' (' + esc(job.location) + ')' : '';
      const durBits = [];
      if (job.period) durBits.push(esc(job.period));
      if (job.duration) durBits.push(esc(ui.labelDuration) + ' : ' + esc(job.duration));
      const right = durBits.length ? ' | ' + durBits.join(' - ') : '';
      const companyHtml = linkedText(job.company, job.website, 'company-link');
      headBits.push('<div class="company-duration">' + companyHtml + locPart + right + '</div>');
    } else if (job.period) {
      headBits.push('<div class="company-duration">' + esc(job.period) + '</div>');
    }

    if (job.teamSize) {
      headBits.push('<div class="team-info">' + esc(ui.labelTeamSize) + ' ' + esc(job.teamSize) + '</div>');
    }
    if (job.project) {
      headBits.push('<div class="project-info">' + esc(ui.labelProject) + ' : ' + esc(job.project) + '</div>');
    }
    if (job.methodology) {
      headBits.push('<div class="project-info">' + esc(ui.labelMethodology) + ' : ' + esc(job.methodology) + '</div>');
    }
    if (job.tags && job.tags.length) {
      headBits.push('<div class="project-info">' + job.tags.map(esc).join(' · ') + '</div>');
    }

    const logo = logoImg(job.website, 'job-logo');
    const parts = [
      '<div class="job-header">' +
        logo +
        '<div class="job-header-text">' + headBits.join('') + '</div>' +
      '</div>'
    ];

    if (job.objective) {
      parts.push(
        '<div class="job-objective"><h4>' + esc(ui.labelObjective) + '</h4><p>' +
        esc(job.objective) + '</p></div>'
      );
    }
    if (job.actions && job.actions.length) {
      const lis = job.actions.map(function (a) { return '<li>' + esc(a) + '</li>'; }).join('');
      parts.push(
        '<div class="job-actions"><h4>' + esc(ui.labelActions) + '</h4><ul>' + lis + '</ul></div>'
      );
    }
    if (job.techEnv && Array.isArray(job.techEnv) && job.techEnv.length) {
      const ps = job.techEnv.map(function (t) {
        return '<p><strong>' + esc(t.category) + ' :</strong> ' + esc(t.value) + '</p>';
      }).join('');
      parts.push(
        '<div class="tech-env"><h4>' + esc(ui.labelTechEnv) + '</h4>' + ps + '</div>'
      );
    } else if (typeof job.techEnv === 'string' && job.techEnv) {
      parts.push(
        '<div class="tech-env"><h4>' + esc(ui.labelTechEnv) + '</h4><p>' + esc(job.techEnv) + '</p></div>'
      );
    }

    return '<div class="job">' + parts.join('') + '</div>';
  }

  function renderFreelanceProject(fp, ui) {
    return [
      '<div class="freelance-project">',
      '  <div class="project-title">' + esc(fp.title) + '</div>',
      '  <div class="client-duration">' +
            esc(ui.labelClient) + ' : ' + esc(fp.client) +
            ' - ' + esc(ui.labelDuration) + ' : ' + esc(fp.duration) +
          '</div>',
      '  <p><strong>' + esc(ui.labelObjective) + ' :</strong> ' + esc(fp.objective) + '</p>',
      '  <p><strong>' + esc(ui.labelTechEnv) + ' :</strong> ' + esc(fp.techEnv) + '</p>',
      '</div>'
    ].join('\n');
  }

  function renderExtended(data) {
    const meta = data.meta;
    const ui = data.ui;

    const eduBlock = data.education.map(function (e) {
      const schoolHtml = linkedText(e.school, e.website, 'school-link');
      return [
        '<div class="education-item">',
        '  <div class="education-date">' + esc(e.date) + '</div>',
        '  <div class="education-degree">' + esc(e.degree) + '</div>',
        '  <div class="education-school">' + schoolHtml + '</div>',
        '</div>'
      ].join('\n');
    }).join('');

    const skillRows = data.skills.map(function (s) {
      return [
        '<div class="skills-grid">',
        '  <div class="skill-category">' + esc(s.category) + '</div>',
        '  <div class="skill-list">' + esc(s.items) + '</div>',
        '</div>'
      ].join('\n');
    }).join('');

    const funcSkills = data.functionalSkills.map(function (f) {
      return '<li>' + esc(f) + '</li>';
    }).join('');

    const jobsHtml = data.experience.map(function (j) { return renderJob(j, ui); }).join('');

    const fpHtml = data.freelanceProjects.map(function (fp) {
      return renderFreelanceProject(fp, ui);
    }).join('');

    const overallTech = data.freelanceOverallTech.map(function (t) {
      return '<p><strong>' + esc(t.category) + ' :</strong> ' + esc(t.value) + '</p>';
    }).join('');

    const jr = renderJob(data.juniorRole, ui);

    return [
      renderHeader(meta),
      '<div class="section">',
      '  <h2 class="section-title">' + esc(ui.sectionEducation) + '</h2>',
      '  ' + eduBlock,
      '</div>',
      '<div class="section">',
      '  <h2 class="section-title">' + esc(ui.sectionTechSkills) + '</h2>',
      '  ' + skillRows,
      '</div>',
      '<div class="section">',
      '  <h2 class="section-title">' + esc(ui.sectionFunctionalSkills) + '</h2>',
      '  <ul class="functional-skills">' + funcSkills + '</ul>',
      '</div>',
      '<div class="section">',
      '  <h2 class="section-title">' + esc(ui.sectionExperience) + '</h2>',
      '  ' + jobsHtml,
      '</div>',
      '<div class="section">',
      '  <h3 class="freelance-section-title">' + esc(ui.sectionFreelance) + ' :</h3>',
      '  ' + fpHtml,
      '  <div class="tech-env">',
      '    <h4>' + esc(ui.sectionOverallTech) + '</h4>',
      '    ' + overallTech,
      '  </div>',
      '</div>',
      '<div class="section">',
      '  ' + jr,
      '</div>'
    ].join('\n');
  }

  // ---------- Banner state ----------

  function updateBanner(data, params) {
    const fr = document.getElementById('lang-fr');
    const en = document.getElementById('lang-en');
    const br = document.getElementById('view-brief');
    const ext = document.getElementById('view-ext');

    [fr, en, br, ext].forEach(function (el) { if (el) el.classList.remove('active'); });
    if (params.lang === 'fr' && fr) fr.classList.add('active');
    if (params.lang === 'en' && en) en.classList.add('active');
    if (params.view === 'brief' && br) br.classList.add('active');
    if (params.view === 'extended' && ext) ext.classList.add('active');

    if (fr) fr.setAttribute('aria-current', params.lang === 'fr' ? 'page' : 'false');
    if (en) en.setAttribute('aria-current', params.lang === 'en' ? 'page' : 'false');
    if (br) br.setAttribute('aria-current', params.view === 'brief' ? 'page' : 'false');
    if (ext) ext.setAttribute('aria-current', params.view === 'extended' ? 'page' : 'false');

    if (br) br.textContent = data.ui.viewBrief;
    if (ext) ext.textContent = data.ui.viewExtended;
    const printLabel = document.querySelector('#print-btn .print-label');
    if (printLabel) printLabel.textContent = data.ui.printLabel;
  }

  async function paint() {
    const params = parseParams();
    const root = document.getElementById('cv-root');
    root.innerHTML = '<div class="cv-loading">…</div>';
    try {
      const data = await loadData(params.lang);
      document.documentElement.lang = params.lang;
      document.title = data.ui.docTitle;
      document.body.classList.toggle('cv-brief', params.view === 'brief');
      document.body.classList.toggle('cv-extended', params.view === 'extended');
      root.innerHTML = params.view === 'brief' ? renderBrief(data) : renderExtended(data);
      updateBanner(data, params);
      track('CV Navigate', { props: { lang: params.lang, view: params.view } });
    } catch (err) {
      console.error(err);
      root.innerHTML = '<div class="cv-error">Failed to load CV data.</div>';
    }
  }

  function wireBanner() {
    const changes = {
      'lang-fr':   { lang: 'fr' },
      'lang-en':   { lang: 'en' },
      'view-brief': { view: 'brief' },
      'view-ext':   { view: 'extended' }
    };
    Object.keys(changes).forEach(function (id) {
      const el = document.getElementById(id);
      if (!el) return;
      el.addEventListener('click', function (e) {
        e.preventDefault();
        const cur = parseParams();
        const next = Object.assign({}, cur, changes[id]);
        if (next.lang === cur.lang && next.view === cur.view) return;
        setUrl(next.lang, next.view);
        paint();
      });
    });
    const printBtn = document.getElementById('print-btn');
    if (printBtn) printBtn.addEventListener('click', function () {
      const p = parseParams();
      track('CV Print', { props: { lang: p.lang, view: p.view } });
      window.print();
    });
    window.addEventListener('popstate', paint);
  }

  document.addEventListener('DOMContentLoaded', function () {
    wireBanner();
    paint();
  });
})();

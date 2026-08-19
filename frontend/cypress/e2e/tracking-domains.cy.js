const pageURL = '/admin/settings/tracking-domains';
const domainsURL = `${pageURL}?tab=domains`;

const baseDomain = {
  id: 1,
  domain: 'example.com',
  base_domain: 'example.com',
  status: 'pending',
  dns_record_type: 'CNAME',
  dns_record_name: '@',
  dns_record_value: 'tracking.myplatform.com',
  verified_at: null,
  last_error: '',
  created_at: '2026-01-01T10:00:00Z',
  updated_at: '2026-01-01T10:00:00Z',
};

const pendingDomain = {
  id: 2,
  domain: 'click.example.com',
  base_domain: 'example.com',
  status: 'pending',
  dns_record_type: 'CNAME',
  dns_record_name: 'click',
  dns_record_value: 'tracking.myplatform.com',
  verified_at: null,
  last_error: '',
  created_at: '2026-01-01T10:00:00Z',
  updated_at: '2026-01-01T10:00:00Z',
};

const verifiedDomain = {
  ...pendingDomain,
  status: 'verified',
  verified_at: '2026-01-02T10:00:00Z',
};

let domains = [];

const mockAPI = () => {
  cy.intercept('GET', '/api/settings*', {
    body: { data: { 'app.from_email': 'Nexuses <noreply@nexuses.in>' } },
  }).as('getSettings');

  cy.intercept('GET', '/api/tracking-domains*', (req) => {
    req.reply({ body: { data: domains } });
  }).as('getDomains');

  cy.intercept('POST', '/api/tracking-domains', (req) => {
    if (req.body.tracking_host) {
      expect(req.body.domain).to.equal('example.com');
      expect(req.body.tracking_host).to.equal(pendingDomain.domain);
      domains = [pendingDomain];
      req.reply({ body: { data: pendingDomain } });
      return;
    }
    expect(req.body.domain).to.equal('example.com');
    domains = [baseDomain];
    req.reply({ body: { data: baseDomain } });
  }).as('addDomain');

  cy.intercept('POST', '/api/tracking-domains/1/verify', (req) => {
    expect(req.body.domain).to.equal('emailtrack.example.com');
    const verified = {
      ...baseDomain,
      domain: 'emailtrack.example.com',
      dns_record_name: 'emailtrack',
      status: 'verified',
      verified_at: '2026-01-02T10:00:00Z',
    };
    domains = [verified];
    req.reply({ body: { data: { ...verified, verified: true, message: 'Domain verified.' } } });
  }).as('verifyBaseDomain');

  cy.intercept('POST', '/api/tracking-domains/2/verify', (req) => {
    expect(req.body.domain).to.equal(pendingDomain.domain);
    domains = [verifiedDomain];
    req.reply({ body: { data: { ...verifiedDomain, verified: true, message: 'Domain verified.' } } });
  }).as('verifyDomain');

  cy.intercept('DELETE', '/api/tracking-domains/*', (req) => {
    domains = [];
    req.reply({ body: { data: true } });
  }).as('deleteDomain');
};

describe('Senders and domains', () => {
  it('Opens the senders and domains page', () => {
    cy.resetDB();
    domains = [];
    mockAPI();
    cy.loginAndVisit(pageURL);

    cy.get('.tracking-domains-brevo__title').should('contain', 'Senders & Domains');
    cy.get('[data-cy=tab-senders]').should('have.class', 'is-active');
    cy.get('[data-cy=tab-domains]').should('exist');
    cy.get('.menu a[data-cy=tracking-domains]').should('contain', 'Senders & Domains');
  });

  it('Opens the add sender page', () => {
    mockAPI();
    cy.loginAndVisit(pageURL);
    cy.get('[data-cy=btn-header-add]').click();
    cy.get('.sender-add__title').should('contain', 'Add Sender');
    cy.get('[data-cy=from-name]').should('exist');
    cy.get('[data-cy=from-email]').should('exist');
    cy.get('.sdd-wizard__cancel').should('contain', 'Cancel');
  });

  it('Adds a domain with branded subdomain', () => {
    mockAPI();
    cy.loginAndVisit(domainsURL);
    cy.wait('@getDomains');

    cy.get('[data-cy=btn-header-add]').click();
    cy.get('[data-cy=domain]').type('example.com');
    cy.get('[data-cy=btn-continue]').click();
    cy.get('[data-cy=subdomain]').type('click');
    cy.get('[data-cy=btn-continue]').click();
    cy.get('[data-cy=btn-continue]').click();
    cy.get('[data-cy=btn-add]').click();
    cy.wait('@addDomain');
    cy.wait('@getDomains');

    cy.get('[data-cy=domain-name]').should('contain', pendingDomain.domain);
    cy.get('[data-cy=status]').should('contain', 'Not authenticated');
  });

  it('Adds a base domain without branded subdomain', () => {
    mockAPI();
    cy.loginAndVisit(domainsURL);
    cy.wait('@getDomains');

    cy.get('[data-cy=btn-header-add]').click();
    cy.get('[data-cy=domain]').type('example.com');
    cy.get('[data-cy=btn-continue]').click();
    cy.get('[data-cy=btn-skip-subdomain]').click();
    cy.wait('@addDomain');
    cy.wait('@getDomains');

    cy.get('[data-cy=domain-name]').should('contain', 'example.com');
    cy.get('[data-cy=status]').should('contain', 'Not authenticated');
  });

  it('Verifies a base domain with tracking host at authentication', () => {
    domains = [baseDomain];
    mockAPI();
    cy.loginAndVisit(domainsURL);
    cy.wait('@getDomains');

    cy.get('[data-cy=btn-verify]').click();
    cy.get('[data-cy=tracking-host]').type('emailtrack.example.com');
    cy.get('.btn-new').contains('Verify DNS').click();
    cy.wait('@verifyBaseDomain');

    cy.get('[data-cy=status]').should('contain', 'Authenticated');
  });

  it('Verifies a domain with preset tracking host', () => {
    domains = [pendingDomain];
    mockAPI();
    cy.loginAndVisit(domainsURL);
    cy.wait('@getDomains');

    cy.get('[data-cy=btn-verify]').click();
    cy.get('.btn-new').contains('Verify DNS').click();
    cy.wait('@verifyDomain');

    cy.get('[data-cy=status]').should('contain', 'Authenticated');
  });

  it('Removes a domain', () => {
    domains = [verifiedDomain];
    mockAPI();
    cy.loginAndVisit(domainsURL);
    cy.wait('@getDomains');

    cy.get('[data-cy=btn-delete]').click();
    cy.get('.modal button.is-primary').click();
    cy.wait('@deleteDomain');
    cy.wait('@getDomains');

    cy.get('[data-cy=domain-card]').should('not.exist');
    cy.get('.bv-empty').should('exist');
  });
});

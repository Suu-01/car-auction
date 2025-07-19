describe('자동차 경매 전체 플로우', () => {
  const id = Date.now()
  const seller = { email: `seller+${id}@test.com`, password: 'pw1234!' }
  const bidder = { email: `bidder+${id}@test.com`, password: 'pw1234!' }
  let auctionId: number

  const api = Cypress.env('apiUrl') || 'http://localhost:8080'

  before(() => {
    // 1) 계정 미리 생성
    cy.request({
      method: 'POST',
      url:    `${api}/users/signup`,     // ← 백틱으로 감싸야 interpolation
      body:   { ...seller, role: 'seller' },
      failOnStatusCode: false,        // 필요시 추가
    })

    cy.request({
      method: 'POST',
      url:    `${api}/users/signup`,     // ← 동일하게 수정
      body:   { ...bidder, role: 'bidder' },
      failOnStatusCode: false,
    })
  })

  it('판매자 로그인 → 출품 페이지 접근', () => {
    cy.login(seller.email, seller.password)
    cy.visit('/auctions/create')

    cy.get('input[placeholder="제목"]').should('be.visible')
    cy.get('textarea[placeholder="설명"]').should('be.visible')
    cy.get('input[placeholder="시작가"]').should('be.visible')
  })

  it('경매 생성 → ID 저장', () => {
    cy.login(seller.email, seller.password)
    cy.visit('/auctions/create')

    const later = new Date(Date.now() + 3600 * 1000)
      .toISOString()
      .slice(0, 16)

    cy.get('input[placeholder="제목"]').type('Cypress Car')
    cy.get('textarea[placeholder="설명"]').type('테스트용 경매')
    cy.get('input[placeholder="시작가"]').type('50000')
    cy.get('input[placeholder="메이커"]').type('Toyota')
    cy.get('input[placeholder="차종명"]').type('Corolla')
    cy.get('input[placeholder="주행거리"]').type('10000')
    cy.get('input[placeholder="연식"]').type('2020')
    cy.get('input[placeholder="사진 URL"]').type('https://via.placeholder.com/150')
    cy.get('input[type=datetime-local]').type(later)

    cy.get('button').contains('출품하기').click()

    // 생성 직후, 응답의 location 헤더나 리다이렉트 ID를 잡아냅니다
    cy.location('pathname').should('match', /\/auctions\/\d+/)
    cy.location('pathname').then((path) => {
      auctionId = Number(path.split('/').pop())
    })
  })

  it('입찰자 로그인 → 목록에서 방금 생성한 경매 확인', () => {
    cy.login(bidder.email, bidder.password)
    cy.visit('/auctions')

    cy.contains('Cypress Car').should('be.visible')
    cy.get('li').contains('Cypress Car').click()
    cy.location('pathname').should('eq', `/auctions/${auctionId}`)
  })

  it('상세 페이지에서 12번 입찰 → 페이지네이션 검증', () => {
    cy.login(bidder.email, bidder.password)
    cy.visit(`/auctions/${auctionId}`)

    // 12번 연속 입찰
    for (let i = 1; i <= 12; i++) {
      cy.get('input[type=number]').clear().type(String(100 + i))
      cy.get('button').contains('입찰하기').click()
      cy.wait(100) // 적당히
    }

    // “다음” 버튼 누르고 2페이지 입찰 목록 확인
    cy.contains('다음').click()
    cy.get('ul > li').first().contains('106') // 100+6
  })
})
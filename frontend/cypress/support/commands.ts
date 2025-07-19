// 1) 타입 보강 (TypeScript 용)
declare namespace Cypress {
  interface Chainable {
    /**
     * @description 커스텀 로그인 커맨드
     * @param email  사용자 이메일
     * @param password 비밀번호
     */
    login(email: string, password: string): Chainable<void>
  }
}

// 2) 커맨드 구현
Cypress.Commands.add('login', (email: string, password: string) => {
  // API_URL 은 cypress.config.cjs 의 env.apiUrl 과 맞춰 주세요.
  const api = Cypress.env('apiUrl') || 'http://localhost:8080'

  // 로그인 API 호출 후 토큰 같은 세션 정보 저장
  cy.request({
    method: 'POST',
    url:    `${api}/users/login`,
    body:   { email, password },
  }).then((resp) => {
    // 예: 응답에 토큰이 body.token 으로 돌아온다면:
    window.localStorage.setItem('accessToken', resp.body.token)
  })
})
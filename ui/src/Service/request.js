function authHeaders(method) {
    if (!method || method === "") {
        method = "GET"
    }
    return {
        method: method,
        headers: {
            'Authorization': 'JWT eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.irAq7TgJOmJu47pesNd6CVn58350r-ntaSvD2YJxohZBIcunGGooNmmGGmp1QOkaQWzsYWJtN8u02wSuT2iABA',
            'Content-Type': 'application/json'
        },
    }
}

export default authHeaders

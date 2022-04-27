const compare = require('./compare')

test('test compare functions', () => {
    const a = {
        spec: {
            date: "2022-03-14T05:37:38Z",
        }
    }
    const b = {
        spec: {
            date: "2022-03-15T05:37:38Z",
        }
    }
    expect(compare.compare(a, b)).toBe(-1)
    expect(compare.compare(b, a)).toBe(1)

    expect([a, b].sort(compare.compare)).toStrictEqual([a, b])
    expect([b, a].sort(compare.compare)).toStrictEqual([a, b])
    expect([a, b].sort(compare.compareRevert)).toStrictEqual([b, a])
})
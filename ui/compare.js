/**
 * Smaller date comes first
 * @param a
 * @param b
 * @returns {number}
 */
function compare(a, b) {
    if (new Date(a.spec.date) < new Date(b.spec.date)){
        return -1;
    } else {
        return 1
    }
}

/**
 * Bigger date comes first
 * @param a
 * @param b
 * @returns {number}
 */
function compareRevert(a, b) {
    return compare(a, b) * -1
}

exports.compare = compare
exports.compareRevert = compareRevert

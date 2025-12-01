import java.io.File

private const val START_ROTATION = 50
private const val MAX_ROTATION = 100

fun main() {
    // 1. DEFINITIONS
    // ---------------------------------------------------------------------
    fun part1(input: List<String>): Int {
        var rotation = START_ROTATION
        var zeroRotations = 0

        input
            .map { line -> line.first() to line.substring(1).toInt() }
            .onEach { (direction, shift) ->
                rotation = when (direction) {
                    'L' -> (rotation - (shift % MAX_ROTATION)).let { if (it < 0) it + MAX_ROTATION else it }
                    'R' -> (rotation + shift) % MAX_ROTATION
                    else -> error("Unknown direction: $direction")
                }

                if (rotation == 0) {
                    zeroRotations++
                }
            }

        return zeroRotations
    }

    fun part2(input: List<String>): Int {
        var rotation = START_ROTATION
        var zeroRotations = 0

        input
            .map { line -> line.first() to line.substring(1).toInt() }
            .onEach { (direction, shift) ->

                when (direction) {
                    'L' -> {
                        zeroRotations += shift / MAX_ROTATION
                        var newRotation = rotation - (shift % MAX_ROTATION)
                        if (newRotation < 0) {
                            newRotation += MAX_ROTATION

                            if (rotation != 0) {
                                zeroRotations++
                            }
                        }
                        if (newRotation == 0) {
                            zeroRotations++
                        }

                        rotation = newRotation
                    }

                    'R' -> {
                        rotation += shift
                        zeroRotations += rotation / MAX_ROTATION
                        rotation %= MAX_ROTATION
                    }

                    else -> error("Unknown direction: $direction")
                }
            }

        return zeroRotations
    }

    // 2. INPUTS
    val testInput = """
        L68
        L30
        R48
        L5
        R60
        L55
        L1
        L99
        R14
        L82
    """.trimIndent().lines()

    val realInput = File("src\\main\\kotlin\\day1\\InputData.txt").readLines()


    // 3. VERIFICATION & EXECUTION
    // ---------------------------------------------------------------------
    
    // Verify Part 1 against the test case (Replace 999 with expected test result)
    check(part1(testInput) == 3) { "Part 1 Test Failed" }
    println("Part 1 Test passed")
    println("Part 1 Result: ${part1(realInput)}")


    println("---")

    // Verify Part 2 against the test case
    check(part2(testInput) == 6) { "Part 2 Test Failed" }
    println("Part 2 Test passed")
    println("Part 2 Result: ${part2(realInput)}")
}
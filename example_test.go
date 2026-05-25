package relax_test

import (
	"errors"
	"fmt"

	"github.com/luckyman42/relax"
)

type userNotFoundError struct {
	ID int
}

func (e *userNotFoundError) Error() string {
	return fmt.Sprintf("user %d not found", e.ID)
}

func fetchUser(id int) (string, error) {
	return "", errors.New("database unavailable")
}

func Example_quickStart() {
	err := relax.CheckFailer(func() {
		relax.FailWith(errors.New("database unavailable"))
	})

	fmt.Println(err)

	// Output:
	// database unavailable
}

func Example_errorsAs() {
	_, err := relax.CheckValue(func() string {
		relax.FailWith(&userNotFoundError{ID: 42})
		return ""
	})

	var target *userNotFoundError
	fmt.Println(errors.As(err, &target))
	fmt.Println(target.ID)

	// Output:
	// true
	// 42
}

func Example_realisticServiceFlow() {
	E := func() error {
		return errors.New("storage unavailable")
	}

	D := func() {
		relax.FailWith(E())
	}

	C := func() { D() }
	B := func() { C() }
	A := func() error { return relax.CheckFailer(B) }

	fmt.Println(A())

	// Output:
	// storage unavailable
}

func ExampleFailOnError() {
	loadProfile := func() (string, error) {
		return "alice", nil
	}

	profile := relax.FailOnError(loadProfile())

	fmt.Println(profile)

	// Output:
	// alice
}

func ExampleCheckValue2() {
	left, right, err := relax.CheckValue2(func() (int, string) {
		return 7, "ok"
	})

	fmt.Println(left)
	fmt.Println(right)
	fmt.Println(err == nil)

	// Output:
	// 7
	// ok
	// true
}

func ExampleCheckValue3() {
	major, minor, patch, err := relax.CheckValue3(func() (int, int, int) {
		return 1, 2, 3
	})

	fmt.Println(major)
	fmt.Println(minor)
	fmt.Println(patch)
	fmt.Println(err == nil)

	// Output:
	// 1
	// 2
	// 3
	// true
}

// Because Failer implements Unwrap, callers can inspect the returned error
// directly with errors.As without knowing about relax.Failer.
func ExampleCheckValue() {
	profile, err := relax.CheckValue(func() string {
		return relax.FailOnError(fetchUser(42))
	})

	fmt.Println(profile == "")
	fmt.Println(err)

	// Output:
	// true
	// database unavailable
}

func ExampleCheckFailer() {
	err := relax.CheckFailer(func() {
		relax.FailWith(errors.New("something failed"))
	})

	fmt.Println(err)

	// Output:
	// something failed
}

func ExampleCheckError() {
	loadUser := func(id int) (string, error) {
		return "", errors.New("user not found")
	}

	err := relax.CheckError(func() error {
		_ = relax.FailOnError(loadUser(99))
		return nil
	})

	fmt.Println(err)

	// Output:
	// user not found
}

func ExampleConvertToFailer() {
	err := errors.New("boom")

	failer := relax.ConvertToFailer(err)

	fmt.Println(failer.Err)

	// Output:
	// boom
}

func ExampleIsFailer() {
	failer := relax.ConvertToFailer(errors.New("failure"))

	fmt.Println(relax.IsFailer(failer))
	fmt.Println(relax.IsFailer(errors.New("normal error")))

	// Output:
	// true
	// false
}

func ExampleFailer() {
	err := relax.CheckFailer(func() {
		failer := relax.ConvertToFailer(errors.New("repository failed"))

		failer.Fail(
			"repository", "users",
			"operation", "find",
		)
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
		fmt.Println(failer.Context["repository"])
	}

	// Output:
	// repository failed
	// users
}

func ExampleFailWith() {
	err := relax.CheckFailer(func() {
		relax.FailWith(errors.New("save failed"))
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
	}

	// Output:
	// save failed
}

func ExampleFailWith_function() {
	badFun := func() error {
		return errors.New("Failer error")
	}

	err := relax.CheckFailer(func() {
		relax.FailWith(badFun())
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
	}

	// Output:
	// Failer error
}

func ExampleFailWith_context() {
	err := relax.CheckFailer(func() {
		relax.FailWith(errors.New("save failed"),
			"user_id", 42,
			"operation", "save_user",
		)
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
		fmt.Println(failer.Context["user_id"])
		fmt.Println(failer.Context["operation"])
	}

	// Output:
	// save failed
	// 42
	// save_user
}

func ExampleFailOnError2() {
	loadName := func() (string, string, error) {
		return "Ada", "Lovelace", nil
	}

	first, last := relax.FailOnError2(loadName())

	fmt.Println(first)
	fmt.Println(last)

	// Output:
	// Ada
	// Lovelace
}

func ExampleFailOnError3() {
	loadVersion := func() (int, int, int, error) {
		return 1, 2, 3, nil
	}

	major, minor, patch := relax.FailOnError3(loadVersion())

	fmt.Println(major)
	fmt.Println(minor)
	fmt.Println(patch)

	// Output:
	// 1
	// 2
	// 3
}

func ExampleCheckResult() {
	value, err := relax.CheckResult(func() (int, error) {
		if true {
			relax.FailWith(errors.New("calculation failed"))
		}

		return 42, nil
	})

	fmt.Println(value)
	fmt.Println(err)

	// Output:
	// 0
	// calculation failed
}

func ExampleCheckResult2() {
	left, right, err := relax.CheckResult2(func() (int, string, error) {
		return 7, "ok", nil
	})

	fmt.Println(left)
	fmt.Println(right)
	fmt.Println(err == nil)

	// Output:
	// 7
	// ok
	// true
}

func ExampleCheckResult3() {
	major, minor, patch, err := relax.CheckResult3(func() (int, int, int, error) {
		return 1, 2, 3, nil
	})

	fmt.Println(major)
	fmt.Println(minor)
	fmt.Println(patch)
	fmt.Println(err == nil)

	// Output:
	// 1
	// 2
	// 3
	// true
}

func ExampleHandleFailer() {
	relax.HandleFailer(func() {
		relax.FailWith(errors.New("worker failed"))
	}, func(err error) {
		fmt.Println(err)
	})

	// Output:
	// worker failed
}

func ExampleHandleFailer_goroutine() {
	done := make(chan struct{})

	go relax.HandleFailer(func() {
		relax.FailWith(errors.New("worker failed"))
	}, func(err error) {
		fmt.Println(err)
		close(done)
	})

	<-done

	// Output:
	// worker failed
}

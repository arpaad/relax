package relax

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestCheckValue_CatchAndRethrowPreservesFailOnErrorError(t *testing.T) {
	_, err := CheckValue(func() string {
		_, innerErr := CheckValue(func() string {
			return FailOnError("ok", errors.New("must fail"))
		})
		if innerErr != nil {
			FailWith(innerErr)
		}
		return ""
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	var nested Failer
	if errors.As(failer.Err, &nested) {
		t.Error("Expected underlying Err not to be a nested Failer")
	}
}

func TestCheckValue_Success(t *testing.T) {
	result, err := CheckValue(func() int {
		return 42
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestCheckValue_Failer(t *testing.T) {
	result, err := CheckValue(func() int {
		FailWith(errors.New("thrown error"))
		return 0
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCheckValue_FailerPointer(t *testing.T) {
	result, err := CheckValue(func() int {
		FailWith(&Failer{Err: errors.New("thrown error")})
		return 0
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCheckValue_PointerToAFailer(t *testing.T) {
	result, err := CheckValue(func() int {
		panic(&Failer{Err: errors.New("thrown error")})
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCheckValue_OtherPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected re-panic, but none occurred")
		}
		if r != "other panic" {
			t.Errorf("Expected 'other panic', got %v", r)
		}
	}()

	_, _ = CheckValue(func() int {
		panic("other panic")
	})
}

func TestCheckValue2_Success(t *testing.T) {
	result1, result2, err := CheckValue2(func() (int, string) {
		return 42, "ok"
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result1 != 42 || result2 != "ok" {
		t.Fatalf("Expected (42, ok), got (%d, %s)", result1, result2)
	}
}

func TestCheckValue3_Success(t *testing.T) {
	result1, result2, result3, err := CheckValue3(func() (int, string, bool) {
		return 42, "ok", true
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result1 != 42 || result2 != "ok" || !result3 {
		t.Fatalf("Expected (42, ok, true), got (%d, %s, %v)", result1, result2, result3)
	}
}

func TestCheckResult_Success(t *testing.T) {
	result, err := CheckResult(func() (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestCheckResult_Failer(t *testing.T) {
	result, err := CheckResult(func() (int, error) {
		FailWith(errors.New("thrown error"))
		return 0, nil
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCheckResult2_Success(t *testing.T) {
	result1, result2, err := CheckResult2(func() (int, string, error) {
		return 42, "ok", nil
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result1 != 42 || result2 != "ok" {
		t.Fatalf("Expected (42, ok), got (%d, %s)", result1, result2)
	}
}

func TestCheckResult2_ReturnedError(t *testing.T) {
	_, _, err := CheckResult2(func() (int, string, error) {
		return 0, "", errors.New("returned error")
	})
	if err == nil || err.Error() != "returned error" {
		t.Fatalf("Expected returned error, got %v", err)
	}
}

func TestCheckResult3_Success(t *testing.T) {
	result1, result2, result3, err := CheckResult3(func() (int, string, bool, error) {
		return 42, "ok", true, nil
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result1 != 42 || result2 != "ok" || !result3 {
		t.Fatalf("Expected (42, ok, true), got (%d, %s, %v)", result1, result2, result3)
	}
}

func TestCheckResult3_Failer(t *testing.T) {
	result1, result2, result3, err := CheckResult3(func() (int, string, bool, error) {
		FailWith(errors.New("thrown error"))
		return 0, "", false, nil
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result1 != 0 || result2 != "" || result3 {
		t.Fatalf("Expected zero values, got (%d, %s, %v)", result1, result2, result3)
	}
}

func TestCheckError_Success(t *testing.T) {
	err := CheckError(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCheckError_Failer(t *testing.T) {
	err := CheckError(func() error {
		FailWith(errors.New("thrown error"))
		return nil
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
}

func TestCheckFailer_Success(t *testing.T) {
	err := CheckFailer(func() {
		// no thrown error
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCheckFailer_Failer(t *testing.T) {
	err := CheckFailer(func() {
		FailWith(errors.New("unwind0 error"))
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "unwind0 error" {
		t.Errorf("Expected 'unwind0 error', got '%s'", failer.Err.Error())
	}
}

func TestHandleFailer_FailerPropagation(t *testing.T) {
	var mu sync.Mutex
	var got error

	HandleFailer(func() {
		FailWith(errors.New("worker failed"))
	}, func(err error) {
		mu.Lock()
		got = err
		mu.Unlock()
	})

	if got == nil {
		t.Fatal("expected error, got nil")
	}

	if got.Error() != "worker failed" {
		t.Fatalf("unexpected error: %v", got)
	}
}

func TestHandleFailer_NilOnErrorPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when onError is nil")
		}
	}()

	HandleFailer(func() {
		// no-op
	}, nil)
}

func TestHandleFailer_GoroutinePropagation(t *testing.T) {
	done := make(chan error, 1)

	go HandleFailer(func() {
		FailWith(errors.New("async failed"))
	}, func(err error) {
		done <- err
	})

	err := <-done

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "async failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleFailer_GoroutineNilOnErrorPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when onError is nil")
		}
	}()

	HandleFailer(func() {
		// should not run
	}, nil)
}

func TestHandleFailer_NonFailerPanic_Repanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected re-panic, but none occurred")
		}
		if r != "non-failer panic" {
			t.Fatalf("expected 'non-failer panic', got %v", r)
		}
	}()

	HandleFailer(func() {
		panic("non-failer panic")
	}, func(err error) {
		t.Error("onError should not be called for non-Failer panics")
	})
}

func TestHandleFailer_ConcurrentMultipleCalls(t *testing.T) {
	const n = 50

	done := make(chan struct{}, n)
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		go HandleFailer(func() {
			if i%2 == 0 {
				FailWith(fmt.Errorf("fail %d", i))
			}
		}, func(err error) {
			errs <- err
			done <- struct{}{}
		})
	}

	count := 0
	for count < n/2 {
		<-done
		count++
	}

	close(errs)

	for err := range errs {
		if err == nil {
			t.Fatal("unexpected nil error")
		}
	}
}

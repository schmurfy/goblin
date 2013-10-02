package goblin

import (
	"testing"
	"reflect"
        "time"
)

type FakeReporter struct {
	describes []string
	fails []string
	passes []string
	ends int
        executionTime time.Duration
        totalExecutionTime time.Duration
    beginFlag, endFlag bool
}

func (r *FakeReporter) beginDescribe(name string) {
	r.describes = append(r.describes, name)
}

func (r *FakeReporter) endDescribe() {
	r.ends++
}

func (r *FakeReporter) itFailed(name string) {
	r.fails = append(r.fails, name)
}

func (r *FakeReporter) itPassed(name string) {
	r.passes = append(r.passes, name)
}

func (r *FakeReporter) itTook(duration time.Duration) {
    r.executionTime = duration
    r.totalExecutionTime += duration
}

func (r *FakeReporter) begin() {
    r.beginFlag = true
}

func (r *FakeReporter) end() {
    r.endFlag = true
}

func TestReporting(t *testing.T) {
	fakeTest := &testing.T{}
	reporter := FakeReporter{}
	fakeReporter := Reporter(&reporter)

	g := Goblin(fakeTest)
	g.SetReporter(fakeReporter)

	g.Describe("One", func() {
		g.It("Foo", func() {
			g.Assert(0).Equals(1)
		})
		g.Describe("Two", func() {
			g.It("Bar", func() {
				g.Assert(0).Equals(0)
			})
		})
	})


	if !reflect.DeepEqual(reporter.describes, []string{"One", "Two"}) {
		t.FailNow()
	}
	if !reflect.DeepEqual(reporter.fails, []string{"Foo"}) {
		t.FailNow()
	}
	if !reflect.DeepEqual(reporter.passes, []string{"Bar"}) {
		t.FailNow()
	}
	if reporter.ends != 2 {
		t.FailNow()
	}

    if !reporter.beginFlag || !reporter.endFlag {
      t.FailNow()
    }
}


func TestReportingTime(t *testing.T) {
	fakeTest := &testing.T{}
	reporter := FakeReporter{}
	fakeReporter := Reporter(&reporter)

	g := Goblin(fakeTest)
	g.SetReporter(fakeReporter)

	g.Describe("One", func() {
            g.AfterEach(func() {
                //TODO: Make this an assertion
                if int64(reporter.executionTime / time.Millisecond) < 5 || int64(reporter.executionTime / time.Millisecond) >= 6 {
                    t.FailNow()
                }
            })
            g.It("Foo", func() {
                time.Sleep(5 * time.Millisecond)
            })
            g.Describe("Two", func() {
                g.It("Bar", func() {
                    time.Sleep(5 * time.Millisecond)
                })
            })
	})

        if int64(reporter.totalExecutionTime / time.Millisecond) < 10 || int64(reporter.totalExecutionTime / time.Millisecond) >= 11 {
            t.FailNow()
        }
}

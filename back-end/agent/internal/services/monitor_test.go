package services

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("Monitor Service", func() {
	var mon *Monitor

	const (
		agentID   = 777
		agentName = "testName"
		taskID    = 777
	)

	BeforeEach(func() {
		mon = NewMonitor(agentID, agentName)
	})

	Context("GetInfo", func() {
		It("should return info", func() {
			res := mon.GetInfo()
			Expect(res).ToNot(BeNil())
			Expect(res.Name).To(Equal(agentName))
			Expect(res.AgentID).To(Equal(int32(agentID)))
		})
	})
	Context("CompleteTask", func() {
		It("should complete work and increment count of works", func() {
			tasksCount := 10
			for i := 0; i < tasksCount; i++ {
				mon.AddTask(uint(i))
				err := mon.CompleteTask(uint(i))
				Expect(err).ToNot(HaveOccurred())
			}
			mon.AddTask(taskID)
			err := mon.CompleteTask(taskID)
			Expect(err).ToNot(HaveOccurred())
			info := mon.GetInfo()
			Expect(info.LastTaskID).To(Equal(uint(taskID)))
			Expect(info.CompletedTasks).To(Equal(uint(tasksCount + 1)))
		})
		It("should return error task not found", func() {
			err := mon.CompleteTask(taskID)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("Concurrent methods execute", func() {
		It("should concurrently return correct info", func() {
			tasksCount := 10
			wg := sync.WaitGroup{}
			wg.Add(tasksCount)
			for i := 1; i < tasksCount+1; i++ {
				i := i
				go func() {
					defer wg.Done()
					mon.AddTask(uint(i))
					time.Sleep(time.Duration(tasksCount-i+1) * time.Millisecond)
					mon.CompleteTask(uint(i))
				}()
			}
			info := mon.GetInfo()
			Expect(info.CompletedTasks).To(Equal(uint(0)))
			time.Sleep(time.Millisecond * 5)
			info = mon.GetInfo()
			Expect(info.CompletedTasks).To(Equal(uint(5)))
		})
	})
})

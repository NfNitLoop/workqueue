// Package workqueue simplifies running a set of jobs at a bounded level of
// concurrency.
//
// For example: if you want to crawl hunreds of web pages, but want to limit
// your maximum concurrency, you can submit all of your jobs to workqueue.Run()
// and it will handle the concurrency for you.
package workqueue
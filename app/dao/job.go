// ============================================================================
// This is auto-generated by gf cli tool only once. Fill this file as you wish.
// ============================================================================

package dao

import (
	"NoticeServices/app/dao/internal"
)

// jobDao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type jobDao struct {
	*internal.JobDao
}

var (
	// Job is globally public accessible object for table job operations.
	Job = &jobDao{
		internal.Job,
	}
)

// Fill with you ideas below.

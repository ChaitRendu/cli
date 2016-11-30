package v2_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/v2"
	"code.cloudfoundry.org/cli/command/v2/v2fakes"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("CreateUser Command", func() {
	var (
		cmd        v2.CreateUserCommand
		fakeUI     *ui.UI
		fakeConfig *commandfakes.FakeConfig
		fakeActor  *v2fakes.FakeCreateUserActor
		executeErr error
	)

	BeforeEach(func() {
		out := NewBuffer()
		fakeUI = ui.NewTestUI(nil, out, out)
		fakeConfig = new(commandfakes.FakeConfig)
		fakeActor = new(v2fakes.FakeCreateUserActor)

		cmd = v2.CreateUserCommand{
			UI:     fakeUI,
			Config: fakeConfig,
			Actor:  fakeActor,
		}

		cmd.RequiredArgs.Username = "some-user"
		cmd.RequiredArgs.Password = "some-password"

		fakeConfig.ExperimentalReturns(true)
		fakeConfig.BinaryNameReturns("faceman")
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	It("Displays the experimental warning message", func() {
		Expect(fakeUI.Out).To(Say(command.ExperimentalWarning))
	})

	Context("when checking that the user is logged in", func() {
		BeforeEach(func() {
			fakeConfig.AccessTokenReturns("")
			fakeConfig.RefreshTokenReturns("")
		})

		It("returns an error if the check fails", func() {
			Expect(executeErr).To(MatchError(command.NotLoggedInError{
				BinaryName: "faceman",
			}))
		})
	})

	Context("when the user is logged in", func() {
		BeforeEach(func() {
			fakeConfig.AccessTokenReturns("some-access-token")
			fakeConfig.RefreshTokenReturns("some-refresh-token")
		})

		It("creates the user", func() {
			Expect(fakeUI.Out).To(Say(`
Creating user some-user...
OK

TIP: Assign roles with 'faceman set-org-role' and 'faceman set-space-role'.`))
		})

		Context("when an error occurs", func() {
			var returnedErr error

			BeforeEach(func() {
				returnedErr = errors.New("non-translatable error")
				fakeActor.CreateUserReturns(
					v2action.Warnings{
						"warning-1",
						"warning-2",
					},
					returnedErr,
				)
			})

			It("returns the error and all warnings", func() {
				Expect(executeErr).To(MatchError(returnedErr))
				Expect(fakeUI.Err).To(Say("warning-1"))
				Expect(fakeUI.Err).To(Say("warning-2"))
			})

			// TODO: figure out which error messages need to be translated
			PContext("when the error is translatable", func() {
				It("returns a translatable error", func() {})
			})
		})
	})
})

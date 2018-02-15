package v3action_test

import (
	. "code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/actor/v3action/v3actionfakes"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
)

var _ = Describe("Application Manifest Actions", func() {
	var (
		actor                     *Actor
		fakeCloudControllerClient *v3actionfakes.FakeCloudControllerClient
	)

	BeforeEach(func() {
		fakeCloudControllerClient = new(v3actionfakes.FakeCloudControllerClient)
		actor = NewActor(fakeCloudControllerClient, nil, nil, nil)
	})

	Describe("ApplyApplicationManifest", func() {
		Context("when the app exists", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetApplicationsReturns(
					[]ccv3.Application{{GUID: "some-app-guid"}},
					ccv3.Warnings{"some-app-warning"},
					nil)
				fakeCloudControllerClient.CreateApplicationActionsApplyManifestByApplicationReturns(
					"some-job-url",
					ccv3.Warnings{"some-apply-manifest-warning"},
					nil,
				)
				fakeCloudControllerClient.PollJobReturns(
					ccv3.Warnings{"some-poll-job-warning"},
					nil)
			})

			It("uploads the app manifest", func() {
				tempFile, err := ioutil.TempFile("", "manifest")
				Expect(err).ToNot(HaveOccurred())
				_, err = tempFile.Write([]byte(`---
applications:
- name: some-app-name
  instances: 2
`))
				Expect(err).ToNot(HaveOccurred())
				err = tempFile.Close()
				Expect(err).ToNot(HaveOccurred())

				warnings, err := actor.ApplyApplicationManifest(tempFile.Name(), "some-app-guid")
				Expect(err).ToNot(HaveOccurred())
				Expect(warnings).To(Equal(Warnings{"some-app-warning", "some-apply-manifest-warning", "some-poll-job-warning"}))

				Expect(fakeCloudControllerClient.CreateApplicationActionsApplyManifestByApplicationCallCount()).To(Equal(1))
				appInCall, guidInCall := fakeCloudControllerClient.CreateApplicationActionsApplyManifestByApplicationArgsForCall(0)
				Expect(appInCall.Name).To(Equal("some-app-name"))
				Expect(guidInCall).To(Equal("some-app-guid"))

				Expect(fakeCloudControllerClient.PollJobCallCount()).To(Equal(1))
			})
		})
	})
})
